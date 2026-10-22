package pkg

import "errors"

var ErrInvalid = errors.New("неверная задача")
var ErrNotFound = errors.New("задача не найдена")
var ErrConflict = errors.New("ID уже занят")
