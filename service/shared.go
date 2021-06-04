package service

import "github.com/muchlist/risa_restfull/dto"

// hanya mereturn unit yang memiliki case atau sedang down.
func filterGeneral(data *dto.GenUnitResponseList) {
	temp := dto.GenUnitResponseList{}
	for _, gen := range *data {
		if gen.CasesSize > 0 {
			temp = append(temp, gen)
		}
	}
	*data = temp
}
