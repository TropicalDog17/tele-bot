package handler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchMarket24h(t *testing.T) {
	data, err := MockFetchData24h()
	if err != nil {
		t.Errorf("Error fetching data: %v", err)
	}
	require.NotNil(t, data)
	fmt.Println(data[0:10])
}

func TestGetBiggestGainer24h(t *testing.T) {
	data, err := MockFetchData24h()
	if err != nil {
		t.Errorf("Error fetching data: %v", err)
	}
	require.NotNil(t, data)
	gainer := GetTopNBiggestGainer(data, 1)
	require.NotNil(t, gainer)
}
