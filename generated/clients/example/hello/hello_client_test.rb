require 'minitest/autorun'
require_relative './hello_client'
require_relative './protos/hello'

include Example::Hello

class TestHelloWorld < Minitest::Test

  def setup
    @client = HelloWorldClient.new
  end

  def test_say_hello
    assert_equal @client.say_hello(SayHelloRequest.new), SayHelloResponse.new
  end

end