#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Docker GPU 管理平台 - 开发文档 Web 服务器
使用方法: python docs_server.py
访问地址: http://localhost:9999
"""

import os
import sys
import re
from pathlib import Path

# 修复 Windows 控制台编码
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8', errors='replace')

try:
    from flask import Flask, render_template_string
    import markdown
except ImportError:
    print("正在安装依赖...")
    os.system(f"{sys.executable} -m pip install flask markdown -q")
    from flask import Flask, render_template_string
    import markdown

app = Flask(__name__)

# 文档目录
DOCS_DIR = Path(__file__).parent / "docs"


def slugify(text):
    """将中文标题转换为合适的 ID"""
    # 移除特殊字符，只保留中文、英文、数字、横线和下划线
    slug = re.sub(r'[^\w\u4e00-\u9fff\-]', '', text)
    slug = re.sub(r'\s+', '-', slug)
    return slug.lower()


def parse_headings(md_content):
    """解析 Markdown 内容，提取标题列表"""
    headings = []
    lines = md_content.split('\n')
    for line in lines:
        # 匹配 h2 标题
        match = re.match(r'^##\s+(.+)$', line)
        if match:
            title = match.group(1).strip()
            # 移除可能的锚点链接
            title = re.sub(r'\[([^\]]+)\]\([^)]+\)', r'\1', title)
            headings.append({
                'title': title,
                'id': slugify(title),
                'level': 2
            })
    return headings


# HTML 模板
HTML_TEMPLATE = """
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Docker GPU 管理平台 - 开发文档</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        html {
            scroll-behavior: smooth;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f7fa;
        }
        .container {
            display: flex;
            min-height: 100vh;
        }
        /* 侧边栏 */
        .sidebar {
            width: 280px;
            background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
            color: #fff;
            position: fixed;
            height: 100vh;
            overflow-y: auto;
            padding: 20px 0;
            box-shadow: 2px 0 10px rgba(0,0,0,0.1);
            z-index: 100;
        }
        .sidebar-header {
            padding: 0 20px 20px;
            border-bottom: 1px solid rgba(255,255,255,0.1);
            margin-bottom: 20px;
        }
        .sidebar-header h1 {
            font-size: 18px;
            color: #4fc3f7;
            margin-bottom: 5px;
        }
        .sidebar-header p {
            font-size: 12px;
            color: #aaa;
        }
        .nav-section {
            padding: 10px 20px;
        }
        .nav-section-title {
            font-size: 11px;
            text-transform: uppercase;
            color: #888;
            margin-bottom: 10px;
            letter-spacing: 1px;
        }
        .nav-item {
            display: block;
            padding: 10px 15px;
            color: #ccc;
            text-decoration: none;
            border-radius: 6px;
            margin-bottom: 4px;
            font-size: 14px;
            transition: all 0.2s;
            cursor: pointer;
        }
        .nav-item:hover {
            background: rgba(79, 195, 247, 0.1);
            color: #4fc3f7;
        }
        .nav-item.active {
            background: linear-gradient(90deg, #4fc3f7 0%, #29b6f6 100%);
            color: #1a1a2e;
            font-weight: 500;
        }
        /* 主内容区 */
        .main-content {
            margin-left: 280px;
            padding: 40px;
            flex: 1;
            max-width: 1200px;
        }
        .content-wrapper {
            background: #fff;
            border-radius: 12px;
            padding: 50px;
            box-shadow: 0 2px 12px rgba(0,0,0,0.08);
        }
        /* Markdown 样式 */
        .markdown-body h1 {
            font-size: 32px;
            border-bottom: 2px solid #4fc3f7;
            padding-bottom: 15px;
            margin-bottom: 30px;
            color: #1a1a2e;
        }
        .markdown-body h2 {
            font-size: 24px;
            margin-top: 40px;
            margin-bottom: 20px;
            color: #16213e;
            border-left: 4px solid #4fc3f7;
            padding-left: 15px;
            padding-top: 60px;
            margin-top: -60px;
        }
        .markdown-body h3 {
            font-size: 18px;
            margin-top: 30px;
            margin-bottom: 15px;
            color: #333;
        }
        .markdown-body h4 {
            font-size: 16px;
            margin-top: 25px;
            margin-bottom: 12px;
            color: #555;
        }
        .markdown-body p {
            margin-bottom: 16px;
            color: #444;
        }
        .markdown-body table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            font-size: 14px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            border-radius: 8px;
            overflow: hidden;
        }
        .markdown-body th {
            background: linear-gradient(90deg, #1a1a2e 0%, #16213e 100%);
            color: #fff;
            padding: 14px 16px;
            text-align: left;
            font-weight: 500;
        }
        .markdown-body td {
            padding: 12px 16px;
            border-bottom: 1px solid #eee;
        }
        .markdown-body tr:nth-child(even) {
            background: #f8f9fa;
        }
        .markdown-body tr:hover {
            background: #e3f2fd;
        }
        .markdown-body code {
            background: #f4f4f5;
            padding: 2px 8px;
            border-radius: 4px;
            font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
            font-size: 13px;
            color: #e53935;
        }
        .markdown-body pre {
            background: #1a1a2e;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            margin: 20px 0;
        }
        .markdown-body pre code {
            background: transparent;
            color: #a5d6ff;
            padding: 0;
        }
        .markdown-body ul, .markdown-body ol {
            margin: 16px 0;
            padding-left: 30px;
        }
        .markdown-body li {
            margin-bottom: 8px;
        }
        .markdown-body blockquote {
            border-left: 4px solid #4fc3f7;
            padding: 15px 20px;
            margin: 20px 0;
            background: #e3f2fd;
            border-radius: 0 8px 8px 0;
            color: #1565c0;
        }
        .markdown-body hr {
            border: none;
            height: 2px;
            background: linear-gradient(90deg, transparent, #4fc3f7, transparent);
            margin: 40px 0;
        }
        .markdown-body a {
            color: #1976d2;
            text-decoration: none;
            border-bottom: 1px dashed #1976d2;
        }
        .markdown-body a:hover {
            color: #4fc3f7;
            border-bottom-color: #4fc3f7;
        }
        /* 响应式 */
        @media (max-width: 768px) {
            .sidebar {
                width: 100%;
                height: auto;
                position: relative;
            }
            .main-content {
                margin-left: 0;
            }
            .container {
                flex-direction: column;
            }
        }
        /* 滚动条样式 */
        ::-webkit-scrollbar {
            width: 8px;
        }
        ::-webkit-scrollbar-track {
            background: #f1f1f1;
        }
        ::-webkit-scrollbar-thumb {
            background: #4fc3f7;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            <div class="sidebar-header">
                <h1>Docker GPU 管理</h1>
                <p>开发文档 v1.0</p>
            </div>
            <nav class="nav-section">
                <div class="nav-section-title">目录导航</div>
                {% for heading in headings %}
                <a href="#{{ heading.id }}" class="nav-item" data-id="{{ heading.id }}">{{ heading.title }}</a>
                {% endfor %}
            </nav>
            <nav class="nav-section">
                <div class="nav-section-title">快速链接</div>
                <a href="/" class="nav-item active">开发文档</a>
                <a href="http://localhost:8890/swagger/index.html" class="nav-item" target="_blank">Swagger API</a>
            </nav>
        </aside>
        <main class="main-content">
            <div class="content-wrapper">
                <div class="markdown-body">
                    {{ content | safe }}
                </div>
            </div>
        </main>
    </div>
    <script>
        // 为所有 h2 标题添加 ID
        document.querySelectorAll('.markdown-body h2').forEach(function(h2) {
            var text = h2.textContent.trim();
            var id = text.replace(/[^\\w\\u4e00-\\u9fff\\-]/g, '').toLowerCase();
            h2.id = id;
        });

        // 导航点击跳转
        document.querySelectorAll('.nav-item[data-id]').forEach(function(item) {
            item.addEventListener('click', function(e) {
                e.preventDefault();
                var targetId = this.getAttribute('data-id');
                var target = document.getElementById(targetId);
                if (target) {
                    var offset = 20;
                    var targetPosition = target.getBoundingClientRect().top + window.pageYOffset - offset;
                    window.scrollTo({
                        top: targetPosition,
                        behavior: 'smooth'
                    });
                    // 更新高亮
                    document.querySelectorAll('.nav-item').forEach(function(i) {
                        i.classList.remove('active');
                    });
                    this.classList.add('active');
                }
            });
        });

        // 滚动监听，自动高亮当前章节
        var headingIds = [];
        document.querySelectorAll('.markdown-body h2').forEach(function(h2) {
            headingIds.push(h2.id);
        });

        window.addEventListener('scroll', function() {
            var scrollPos = window.scrollY + 100;
            var current = '';

            document.querySelectorAll('.markdown-body h2').forEach(function(h2) {
                if (h2.offsetTop <= scrollPos) {
                    current = h2.id;
                }
            });

            if (current) {
                document.querySelectorAll('.nav-item[data-id]').forEach(function(item) {
                    item.classList.remove('active');
                    if (item.getAttribute('data-id') === current) {
                        item.classList.add('active');
                    }
                });
            }
        });
    </script>
</body>
</html>
"""


def read_markdown_file(filename: str) -> str:
    """读取 Markdown 文件"""
    filepath = DOCS_DIR / filename
    if filepath.exists():
        with open(filepath, 'r', encoding='utf-8') as f:
            return f.read()
    return f"# 文档不存在\n\n文件 `{filename}` 未找到。"


@app.route('/')
def index():
    """主页 - 显示开发文档"""
    md_content = read_markdown_file('DEVELOPMENT.md')

    # 解析标题
    headings = parse_headings(md_content)

    # 转换 Markdown 为 HTML
    html_content = markdown.markdown(
        md_content,
        extensions=[
            'tables',
            'fenced_code',
            'codehilite',
            'nl2br',
            'sane_lists'
        ]
    )

    return render_template_string(HTML_TEMPLATE, content=html_content, headings=headings)


@app.route('/<doc_name>')
def show_doc(doc_name):
    """显示指定文档"""
    if not doc_name.endswith('.md'):
        doc_name += '.md'
    md_content = read_markdown_file(doc_name)
    headings = parse_headings(md_content)
    html_content = markdown.markdown(
        md_content,
        extensions=['tables', 'fenced_code', 'codehilite', 'nl2br']
    )
    return render_template_string(HTML_TEMPLATE, content=html_content, headings=headings)


if __name__ == '__main__':
    # 确保文档目录存在
    DOCS_DIR.mkdir(parents=True, exist_ok=True)

    print("=" * 50)
    print("   Docker GPU 管理平台 - 开发文档服务器")
    print("=" * 50)
    print(f"   文档目录: {DOCS_DIR}")
    print("   访问地址: http://localhost:9999")
    print("   按 Ctrl+C 停止服务")
    print("=" * 50)

    app.run(host='0.0.0.0', port=9999, debug=False)