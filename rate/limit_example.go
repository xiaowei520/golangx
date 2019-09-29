package rate

func Use() {
	test := DefaultLimiter(1, "test")

	//自定义实例, 用来控制 触发限流后 进行阻塞的时间
	_ = NewLimiter(1, "test", 100, 2000)

	//每次业务都进行是否允许操作,如果被限流,会 阻塞 （时间随机）
	test.Allow()

	//自己业务逻辑触发阻塞
	//比如 第三方响应 请求超频,我们可以调用stop 来对我们进行限流
	test.Stop()

	//当第三方响应正常 我们进行恢复
	test.Recover()

}
