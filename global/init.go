package global

var err error

func Init()  error {
	if err := InitK8sClientSet(); err != nil {
		return err
	}
	return nil
}
