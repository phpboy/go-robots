package rpcdemo

import "errors"

type DivRpc struct {
	A,B int
}

func (DivRpc) Div(args DivRpc,result *float64) error {
	if args.B == 0{
		return errors.New("division by 0")
	}
	*result = float64(args.A) / float64(args.B)
	return nil
}