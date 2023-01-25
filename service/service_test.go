package service

import (
	"Ozon/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_CheckUrl(t *testing.T) {
	s := Service{}
	result := models.Result{}
	in := "https://github.com/sirupsen/logrus?ysclid=ldaqpekbsl261059789"
	value := s.CheckUrl(in, &result)
	require.True(t, value)

	in = ""
	value = s.CheckUrl(in, &result)
	require.False(t, value)
}

func TestService_Shorting(t *testing.T) {
	s := Service{}
	maxLength := 0
	arr := s.Shorting()
	for i, _ := range arr {
		maxLength = i
	}
	require.Equal(t, maxLength, 9)
}
