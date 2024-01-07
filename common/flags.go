package common

import (
	"github.com/spf13/cobra"
)

type GetterT func(name string) (interface{}, error)


func handleDefault[T any](value T, err error, _default T) T {
	if err != nil  {
		return _default
	}
	return value
}

func getDefaultValue[T any](_default ...T) T {
	var defaultValue T
	
	if _default != nil {
		return _default[0]
	}
	
	return defaultValue
}

func GetStringFlag(cmd *cobra.Command, name string) (string, error){
	value, err := cmd.Flags().GetString(name)
	return value, err
}

func GetStringFlagOr(cmd *cobra.Command, name string, _default ...string) string {
	value, err := GetStringFlag(cmd, name)
	return handleDefault[string](value, err, getDefaultValue[string](_default...))
}

func GetUInt16Flag(cmd *cobra.Command, name string) (uint16, error) {
	value, err := cmd.Flags().GetUint16(name)
	return value, err
}

func GetUInt16FlagOr(cmd *cobra.Command, name string, _default ...uint16) uint16 {
	value, err := GetUInt16Flag(cmd, name)
	return handleDefault[uint16](value, err, getDefaultValue[uint16](_default...))
}