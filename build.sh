if [ -z $1 ]; then
	echo "Building project..."
	go build -o gooflife.out main.go gamelogic.go 
elif [ "-r" == $1 ]; then 
	echo "Cleaning project..."
	rm *.out
fi



