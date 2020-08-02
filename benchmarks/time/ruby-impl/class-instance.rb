class Customer
		def initialize()
			@name = "Rajika"
			@street = "NYC"
    end
    def getStreet
      puts(@street)
    end
	end
	i = 0
	a = 100
	while i < a do
		i = i + 1
		cust = Customer.new()
		cust.getStreet()
	end
