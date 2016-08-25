require "./protos/hello"

class HelloWorldClient
  def say_hello(request)
    return SayHelloResponse.new
  end
end
