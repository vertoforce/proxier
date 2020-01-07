package proxier

import (
	"context"
	"fmt"
	"testing"
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
