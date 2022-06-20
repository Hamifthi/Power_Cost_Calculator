# Instructions to run app, and a description for decisions

## How to run code
create these directories generated_data, inputs, outputs

unzip inputs to inputs folder to test service with generated data in inputs.zip file

copy the sample.env file to .env file

The go version is 1.17

go mod download && go mod verify

go test -v ./...

go build -v -o ./app ./cmd/main.go

./app

### run with docker
create outputs directory to get generated files in case want to check them.

docker build -t calculator .

docker run -it --rm --name app-container --mount type=bind,source="$(pwd)"/outputs,target=/usr/src/app/outputs calculator

## design decisions
I decided to use clean architecture to write the app because of its extensibility,
 and the well-defined and separated responsibility of each layer.
So I put the models in the entity package, and the business or core logic in the pkg folder
 as services.
Also, I try to use the Go standard project layout from this github repo 
 https://github.com/golang-standards/project-layout

I use dependency injection for each layer dependency, As it makes testing much easier
 and, also is based on dependency inversion from SOLID principles. I try my best to use
 solid principles in the whole code.
 
I put the configurable settings in env file and use a service like code for handling reading
 environment variables because handling errors in the main file makes it dirty.
 It uses logger as dependency.
 
Because the size of sessions may become so large that may not fit in ram totally I use
 buffer and process data in chunks. I use a producer consumer like design where multiple
 chunks processed and wrote down to the output file. By utilizing the sync.pool feature
 I can reduce the pressure on the GC and use less memory. Besides, with the help of channels
 I limit the number of threads that is also configurable to reduce the pressure on CPU.

Since writing down the results may cause race condition between go routines I put a mutex lock
 for output file.

With this process model we can handle files as large as multiple gigabytes of data. The size
 of the ram to allocate is configurable.
  
With the help of the interfaces the implementation become flexible in case in future the
 implementation wants change due to new requirements.

I write unit tests for all parts except for cases that trigger error, where it might not common.

With the help of fake data generator I'm able to test my code in a production like environment,
 with this in mind that fake data generator is just a tool to test the service, I didn't optimize
 the data generator to run faster or consume less memory. It could be done if the project requirements
 change in the future.

## Benchmarks and results
Because the size of sessions could go beyond millions the first approach which was comparing each session
 with each tariff becomes inefficient soon because its time is approximately O(n2).
So I decided to use binary search for finding the tariff that may applied to session and this decrease time
 complexity to approximately O(logn).
I test code with a 300 MB fake sessions that result in runtime of 8 minutes with the first version of code and
 about 2 minute and 30 seconds with the version uses binary search.
