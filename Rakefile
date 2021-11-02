@versions = {}

File.read("version").each_line do |line|
	product,versionnumber = line.chomp.split(/=/)
	@versions[product]=versionnumber
end

@ets_version = @versions['ets_version']

desc "Show rake description"
task :default do
	puts
	puts "Run 'rake -T' for a list of tasks."
	puts
	puts "1: Use 'rake build' to build the 'ets' binary. That should be\n   the starting point."
	puts
end


desc "Compile and install necessary software"
task :build  do
	sh "go build -ldflags \"-X main.version=#{@ets_version}\" -o bin/ets github.com/speedata/ets/ets/ets"
end



