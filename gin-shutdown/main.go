package main

import "sync"

func main() {
	//router := gin.Default()
	//router.GET("/", func(c *gin.Context) {
	//	time.Sleep(5 * time.Second)
	//	c.String(http.StatusOK, "Welcome Gin Server")
	//})
	//
	//srv := &http.Server{
	//	Addr:    ":8080",
	//	Handler: router,
	//}
	//
	//// Initializing the server in a goroutine so that
	//// it won't block the graceful shutdown handling below
	//go func() {
	//	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()
	//
	//// kill (no param) default send syscanll.SIGTERM
	//// kill -2 is syscall.SIGINT
	//// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	//log.Println("Shutdown Server ...")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	//// catching ctx.Done(). timeout of 5 seconds.
	//<-ctx.Done()
	//log.Println("Server exiting")

	group := sync.WaitGroup{}
	group.Add(1)
	group.Wait()
}
