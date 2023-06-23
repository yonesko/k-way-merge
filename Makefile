run:
	go run cmd/uniq_count/main.go -input input -N 20000 && head input_uniq && tail input_uniq && head input_sorted && tail input_sorted