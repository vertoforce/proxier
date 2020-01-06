package getproxylist

import (
	"context"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	p := GetProxyListSource{}
	proxy, err := p.GetProxy(context.Background())
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(proxy)
}
