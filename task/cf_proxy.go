package task

import (
	"CloudflareSpeedTest/utils"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	DisableProxy = false
)

//http clint
func request(method string, url string, body io.Reader) (*http.Response, error) {
	clint := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic("create request error")
	}
	req.Header.Add("Host", "cloudflare.cdn.openbsd.org")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	resp, err := clint.Do(req)
	return resp, err
}
func proxyTest(data utils.CloudflareIPData) bool {
	url := fmt.Sprintf("http://%s/cdn-cgi/trace", data.IP.String())
	if data.IP.IP.To16() != nil {
		url = fmt.Sprintf("http://[%s]/cdn-cgi/trace", data.IP.String())
	}
	r, e := request("GET", url, nil)
	if e != nil || r.StatusCode != 200 {
		return false
	}
	return true
}
func ProxyTestMain(data utils.PingDelaySet) utils.PingDelaySet {
	var result utils.PingDelaySet
	var wg sync.WaitGroup
	var lock sync.Mutex
	ch := make(chan int, Routines)
	if DisableProxy {
		return data
	}
	testNum := len(data)
	br := utils.NewBar(testNum)
	fmt.Printf("开始测试是否为代理ip,数量%d\n", testNum)
	for i := 0; i < testNum; i++ {
		// 在每个 IP 下载测速后，以 [下载速度下限] 条件过滤结果
		ch <- i
		wg.Add(1)
		d := data[i]
		go func(d utils.CloudflareIPData) {
			isSucess := proxyTest(d)
			if isSucess == true {
				lock.Lock()
				result = append(result, d)
				lock.Unlock()
			}
			defer br.Grow(1)
			defer wg.Done()
			<-ch
		}(d)
	}
	wg.Wait()
	br.Done()
	fmt.Printf("可用的ip: %d \n", len(result))
	return result
}
