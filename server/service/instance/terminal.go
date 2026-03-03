package instance

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/computenode"
	instanceModel "github.com/flipped-aurora/gin-vue-admin/server/model/instance"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该限制
	},
}

// TerminalMessage 终端消息
type TerminalMessage struct {
	Type string `json:"type"` // input, resize, ping
	Data string `json:"data"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

// HandleTerminal 处理终端WebSocket连接
func (instanceService *InstanceService) HandleTerminal(c *gin.Context, ID string, shell string) {
	// 升级HTTP连接为WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("WebSocket升级失败", zap.Error(err))
		return
	}
	defer ws.Close()

	// 获取实例和节点信息
	var inst instanceModel.Instance
	if err := global.GVA_DB.Where("id = ?", ID).First(&inst).Error; err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("获取实例信息失败: "+err.Error()))
		return
	}

	if inst.ContainerId == nil || *inst.ContainerId == "" {
		ws.WriteMessage(websocket.TextMessage, []byte("容器ID为空"))
		return
	}

	if inst.NodeId == nil {
		ws.WriteMessage(websocket.TextMessage, []byte("实例未关联节点"))
		return
	}

	var node computenode.ComputeNode
	if err := global.GVA_DB.Where("id = ?", *inst.NodeId).First(&node).Error; err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("获取节点信息失败: "+err.Error()))
		return
	}

	// 创建Docker客户端
	cli, err := dockerService.CreateDockerClient(&node)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("创建Docker客户端失败: "+err.Error()))
		return
	}
	defer cli.Close()

	ctx := context.Background()

	// 根据shell参数确定使用的shell命令
	var shellCmd string
	if shell == "sh" {
		shellCmd = "/bin/sh"
	} else {
		shellCmd = "/bin/bash"
	}

	// 创建exec实例
	execConfig := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{shellCmd},
	}

	// 创建exec，如果失败则尝试另一个shell
	execResp, err := cli.ContainerExecCreate(ctx, *inst.ContainerId, execConfig)
	if err != nil {
		// 如果选择的shell失败，尝试另一个
		if shell == "bash" {
			execConfig.Cmd = []string{"/bin/sh"}
		} else {
			execConfig.Cmd = []string{"/bin/bash"}
		}
		execResp, err = cli.ContainerExecCreate(ctx, *inst.ContainerId, execConfig)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("创建exec失败: "+err.Error()+"\r\n"))
			return
		}
		// 如果fallback成功，发送提示信息
		ws.WriteMessage(websocket.TextMessage, []byte("注意: 容器中不存在 "+shellCmd+"，已切换到备用shell\r\n"))
	}

	// 附加到exec
	attachResp, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{
		Tty: true,
	})
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("附加到exec失败: "+err.Error()))
		return
	}
	defer attachResp.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// 从容器读取输出发送到WebSocket
	go func() {
		defer wg.Done()
		buf := make([]byte, 32*1024) // 增大缓冲区
		for {
			n, err := attachResp.Reader.Read(buf)
			if err != nil {
				// EOF 是正常的流结束标志，不需要记录错误
				global.GVA_LOG.Debug("容器输出流结束", zap.Error(err))
				return
			}
			if n > 0 {
				// 使用 BinaryMessage 发送原始数据
				err = ws.WriteMessage(websocket.BinaryMessage, buf[:n])
				if err != nil {
					global.GVA_LOG.Error("发送WebSocket消息失败", zap.Error(err))
					return
				}
			}
		}
	}()

	// 从WebSocket读取输入发送到容器
	go func() {
		defer wg.Done()
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					global.GVA_LOG.Error("WebSocket读取消息失败", zap.Error(err))
				}
				return
			}

			// 只处理文本消息（JSON格式）
			if messageType == websocket.TextMessage {
				var msg TerminalMessage
				if err := json.Unmarshal(message, &msg); err == nil {
					switch msg.Type {
					case "input":
						// 发送用户输入到容器
						if _, err := attachResp.Conn.Write([]byte(msg.Data)); err != nil {
							global.GVA_LOG.Error("写入容器输入失败", zap.Error(err))
							return
						}
					case "resize":
						// 调整终端大小
						if err := cli.ContainerExecResize(ctx, execResp.ID, container.ResizeOptions{
							Height: uint(msg.Rows),
							Width:  uint(msg.Cols),
						}); err != nil {
							global.GVA_LOG.Error("调整终端大小失败", zap.Error(err))
						}
					case "ping":
						// 心跳检测，发送pong响应
						pongMsg := TerminalMessage{Type: "pong"}
						if pongData, err := json.Marshal(pongMsg); err == nil {
							ws.WriteMessage(websocket.TextMessage, pongData)
						}
					}
				} else {
					// 如果不是有效的JSON，尝试直接作为输入发送
					if _, err := attachResp.Conn.Write(message); err != nil {
						global.GVA_LOG.Error("写入容器输入失败", zap.Error(err))
						return
					}
				}
			}
		}
	}()

	// 设置WebSocket心跳
	ws.SetPingHandler(func(appData string) error {
		return ws.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	wg.Wait()
}
