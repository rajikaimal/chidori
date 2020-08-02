class Customer
		def initialize()
			@name = "Rajika"
			@street = "5th Street"
		end
		def foo
			puts(@name)
			puts(@street)
		end
end
	i = 0
	a = 100
	while i < a do
		i = i + 1
		cust = Customer.new()
		cust.foo()
	end
