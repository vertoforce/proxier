package proxier

import (
	"context"
	"fmt"
	"testing"

	"github.com/vertoforce/proxier/proxy/proxysources/gimmeproxy"
)

const (
	TestURL = "https://www.meetup.com/Huey-Spheres-by-GoHuey-com/events/kjwzvqyzqbgc"
)

func TestDoRequest(t *testing.T) {
	proxier := New()
	resp, err := proxier.DoRequest(context.Background(), "GET", TestURL, nil)
	if err != nil {
		t.Error(err)
		return
	}
	// bytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// ioutil.WriteFile("out.htm", bytes, 0644)
	fmt.Println(resp)
}

func TestCacheProxies(t *testing.T) {
	proxier := New()
	added, err := proxier.CacheProxies(context.Background(), 5)
	if err != nil {
		t.Error(err)
	}

	if added == 0 {
		t.Errorf("Didn't add any proxies")
	}
}

// This test sees how long we can make requests to one of our proxy source URLs
// Normally the URL will start blocking after 10 requests, but that's we use proxies
// It's a bit meta (using proxy to make request to proxy source)
func TestMakeRequestThatGetsDenied(t *testing.T) {
	proxier := New()
	for {
		resp, err := proxier.DoRequest(context.Background(), "GET", gimmeproxy.GimmeProxyURL, nil)
		if err != nil {
			fmt.Printf("Error %s: \n", err.Error())
			break
		}
		fmt.Printf("Response code: %d\n", resp.StatusCode)
	}

}
