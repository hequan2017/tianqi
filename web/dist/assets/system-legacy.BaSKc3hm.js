/*! 
 Build based on tianqi 
 Time : 1772519311000 */
System.register(["./index-legacy.Dd4OSXyr.js"],function(t,e){"use strict";var n;return{setters:[function(t){n=t.s}],execute:function(){t("b",function(){return n({url:"/system/getSystemConfig",method:"post"})}),t("s",function(t){return n({url:"/system/setSystemConfig",method:"post",data:t})}),t("a",function(){return n({url:"/system/getServerInfo",method:"post",donNotShowLoading:!0})}),t("r",function(t){return n({url:"/system/reloadSystem",method:"post",data:t})}),t("g",function(){return n({url:"/system/getJumpboxConfig",method:"post"})})}}});
