/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/

package login

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
)

/**
实例化login module
*/
var Module = func() module.Module {
	login := new(login)
	return login
}

type login struct {
	basemodule.BaseModule
}

func (m *login) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Login"
}
func (m *login) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (m *login) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)

	m.GetServer().RegisterGO("HD_Login", m.login) //我们约定所有对客户端的请求都以Handler_开头
	m.GetServer().RegisterGO("track", m.track)    //演示后台模块间的rpc调用
	m.GetServer().RegisterGO("track2", m.track2)  //演示后台模块间的rpc调用
	m.GetServer().RegisterGO("track3", m.track3)  //演示后台模块间的rpc调用
	m.GetServer().Register("HD_Robot", m.robot)
	m.GetServer().RegisterGO("HD_Robot_GO", m.robot) //我们约定所有对客户端的请求都以Handler_开头
}

func (m *login) Run(closeSig chan bool) {
}

func (m *login) OnDestroy() {
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}
func (m *login) robot(session gate.Session, msg map[string]interface{}) (result string, err string) {
	//time.Sleep(1)
	return "sss", ""
}
func (m *login) login(session gate.Session, msg map[string]interface{}) (result string, err string) {
	if msg["userName"] == nil || msg["passWord"] == nil {
		result = "userName or passWord cannot be nil"
		return
	}
	userName := msg["userName"].(string)
	err = session.Bind(userName)
	if err != "" {
		return
	}
	session.Set("login", "true")
	session.Push() //推送到网关
	return fmt.Sprintf("login success %s", userName), ""
}

func (m *login) track(session gate.Session) (result string, err string) {
	//演示后台模块间的rpc调用
	time.Sleep(time.Millisecond * 10)
	log.TInfo(session, "Login %v", "track1")
	m.RpcInvoke("Login", "track2", session)
	return fmt.Sprintf("My is Login Module "), ""
}

func (m *login) track2(session gate.Session) (result string, err string) {
	//演示后台模块间的rpc调用
	time.Sleep(time.Millisecond * 10)
	log.TInfo(session, "Login %v", "track2")
	r := rand.Intn(100)
	if r > 30 {
		m.RpcInvoke("Login", "track3", session)
	}

	return fmt.Sprintf("My is Login Module"), ""
}
func (m *login) track3(session gate.Session) (result string, err string) {
	//演示后台模块间的rpc调用
	time.Sleep(time.Millisecond * 10)
	log.TInfo(session, "Login %v", "track3")
	return fmt.Sprintf("My is Login Module"), ""
}
