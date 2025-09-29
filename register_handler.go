package main

// func handlerRegister(s *state, cmd command) error {
// 	args := cmd.args
// 	if len(args) != 1 {
// 		return errors.New("need username as argument")
// 	}
// 	name := cmd.args[0]
// 	_, err := s.db.GetUser(context.Background(), name)
// 	if err == nil {
// 		return errors.New("user already exists")
// 	}
//
// 	u, e := s.db.CreateUser(
// 		context.Background(),
// 		database.CreateUserParams{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 			Name:      name,
// 		},
// 	)
// 	if e != nil {
// 		return e
// 	}
// 	s.config.SetUser(u.Name)
// 	fmt.Printf("created user: %v", u)
//
// 	return nil
// }
