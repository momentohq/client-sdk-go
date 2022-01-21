package main

func main() {
	//scsControlEndpoint := "control.cell-alpha-dev.preprod.a.momentohq.com:443"

	// cc, cErr := NewScsControlClient(TEST_AUTH_TOKEN, scsControlEndpoint)
	// if cErr != nil {
	// 	fmt.Println(cErr.Error())
	// }
	// resp, err := cc.ListCaches()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(resp.Caches[0].Name)


	// scsDataEndpoint := "cache.cell-alpha-dev.preprod.a.momentohq.com:443"

	// config := &tls.Config{
	// 	InsecureSkipVerify: false,
	// }
	
	// conn, err := grpc.Dial(scsControlEndpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(addHeadersInterceptor))
	// if err != nil {
	// 	fmt.Println("Something went wrong establishing gRPC channel")
	// 	fmt.Print(err.Error())
	// }

	// connCache, errData := grpc.Dial(scsDataEndpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithDisableRetry(), grpc.WithUnaryInterceptor(addHeadersInterceptor))
	// if errData != nil {
	// 	fmt.Println("Something went wrong establishing gRPC channel")
	// 	fmt.Print(errData.Error())
	// }
	// cacheClient := pb.NewScsClient(conn)
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// controlClient := pb.NewScsControlClient(conn)
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// request := &pb.CreateCacheRequest{CacheName: "cache-go"}
	// resp, er := controlClient.CreateCache(ctx, request)
	// request := &pb.DeleteCacheRequest{CacheName: "cache"}
	// resp, er := controlClient.DeleteCache(ctx, request)
	// request := &pb.ListCachesRequest{}
	// resp, er := controlClient.ListCaches(ctx, request)

	//md := metadata.Pairs("cache", "cache-go")
	// request := &pb.SetRequest{CacheKey: []byte("erika"), CacheBody: []byte("tharp"), TtlMilliseconds: uint32(300000)}
	// resp, er := cacheClient.Set(metadata.NewOutgoingContext(ctx, md), request)
	// getRequest := &pb.GetRequest{CacheKey: []byte("erika")}
	// r, e := cacheClient.Get(metadata.NewOutgoingContext(ctx, md), getRequest)
	// if e != nil {
	// 	fmt.Println("Something went wrong setting cache")
	// 	fmt.Print(e.Error())
	// }

	// fmt.Println("No error detected!")
	// fmt.Println(r)
	
}
