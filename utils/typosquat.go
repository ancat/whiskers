package utils

import (
	"strings"
)

type TyposquatsJSON struct {
	Candidates map[string][]string `json:"possible_typos"`
}

// List of popular gems to check against
var popularGems = []string{
	"Ascii85", "CFPropertyList", "aasm", "actioncable", "actionmailbox",
	"actionmailer", "actionpack", "actiontext", "actionview",
	"active_model_serializers", "activejob", "activemodel", "activerecord",
	"activerecord-import", "activeresource", "activestorage", "activesupport",
	"addressable", "akami", "ansi", "arel", "ast", "atomos", "autoprefixer-rails",
	"awesome_print", "aws-eventstream", "aws-partitions", "aws-sdk", "aws-sdk-acm",
	"aws-sdk-apigateway", "aws-sdk-applicationautoscaling",
	"aws-sdk-applicationdiscoveryservice", "aws-sdk-appstream", "aws-sdk-athena",
	"aws-sdk-autoscaling", "aws-sdk-batch", "aws-sdk-budgets",
	"aws-sdk-cloudformation", "aws-sdk-cloudfront", "aws-sdk-cloudhsm",
	"aws-sdk-cloudhsmv2", "aws-sdk-cloudsearch", "aws-sdk-cloudtrail",
	"aws-sdk-cloudwatch", "aws-sdk-cloudwatchevents", "aws-sdk-cloudwatchlogs",
	"aws-sdk-codebuild", "aws-sdk-codecommit", "aws-sdk-codedeploy",
	"aws-sdk-codepipeline", "aws-sdk-cognitoidentity",
	"aws-sdk-cognitoidentityprovider", "aws-sdk-comprehend",
	"aws-sdk-configservice", "aws-sdk-core", "aws-sdk-costandusagereportservice",
	"aws-sdk-costexplorer", "aws-sdk-databasemigrationservice", "aws-sdk-dax",
	"aws-sdk-devicefarm", "aws-sdk-directconnect", "aws-sdk-directoryservice",
	"aws-sdk-dynamodb", "aws-sdk-dynamodbstreams", "aws-sdk-ec2", "aws-sdk-ecr",
	"aws-sdk-ecs", "aws-sdk-efs", "aws-sdk-elasticache",
	"aws-sdk-elasticbeanstalk", "aws-sdk-elasticloadbalancing",
	"aws-sdk-elasticloadbalancingv2", "aws-sdk-elasticsearchservice",
	"aws-sdk-elastictranscoder", "aws-sdk-emr", "aws-sdk-firehose",
	"aws-sdk-gamelift", "aws-sdk-glacier", "aws-sdk-glue", "aws-sdk-greengrass",
	"aws-sdk-guardduty", "aws-sdk-health", "aws-sdk-iam", "aws-sdk-iot",
	"aws-sdk-iotdataplane", "aws-sdk-kinesis", "aws-sdk-kms", "aws-sdk-lambda",
	"aws-sdk-lex", "aws-sdk-lexmodelbuildingservice", "aws-sdk-lightsail",
	"aws-sdk-marketplacemetering", "aws-sdk-mediaconvert", "aws-sdk-medialive",
	"aws-sdk-migrationhub", "aws-sdk-opsworks", "aws-sdk-opsworkscm",
	"aws-sdk-organizations", "aws-sdk-pinpoint", "aws-sdk-polly", "aws-sdk-rds",
	"aws-sdk-redshift", "aws-sdk-rekognition", "aws-sdk-resourcegroupstaggingapi",
	"aws-sdk-resources", "aws-sdk-route53", "aws-sdk-route53domains", "aws-sdk-s3",
	"aws-sdk-sagemaker", "aws-sdk-secretsmanager", "aws-sdk-servicecatalog",
	"aws-sdk-ses", "aws-sdk-shield", "aws-sdk-simpledb", "aws-sdk-sms",
	"aws-sdk-snowball", "aws-sdk-sns", "aws-sdk-sqs", "aws-sdk-ssm",
	"aws-sdk-states", "aws-sdk-storagegateway", "aws-sdk-swf", "aws-sdk-waf",
	"aws-sdk-workspaces", "aws-sdk-xray", "aws-sigv4", "axiom-types", "babosa",
	"backports", "base64", "bcrypt", "better_errors", "bigdecimal", "bindata",
	"bindex", "binding_of_caller", "bootboot", "bootsnap", "brakeman", "browser",
	"bson", "builder", "bullet", "bundler", "bundler-audit", "byebug",
	"capistrano", "capybara", "capybara-screenshot", "carrierwave", "celluloid",
	"celluloid-essentials", "celluloid-extras", "celluloid-fsm", "celluloid-io",
	"celluloid-pool", "celluloid-supervision", "childprocess", "chronic",
	"chunky_png", "claide", "claide-plugins", "climate_control", "cocoapods",
	"cocoapods-core", "cocoapods-downloader", "coderay", "coercible",
	"coffee-rails", "coffee-script", "coffee-script-source", "colored", "colored2",
	"colorize", "commander", "concurrent-ruby", "connection_pool", "cork",
	"countries", "crack", "crass", "css_parser", "cucumber", "daemons", "dalli",
	"danger", "database_cleaner", "database_cleaner-active_record", "date",
	"ddtrace", "debase-ruby_core_source", "debug_inspector", "declarative",
	"declarative-option", "deep_merge", "descendants_tracker", "devise",
	"diff-lcs", "diffy", "digest-crc", "docile", "dogapi", "dogstatsd-ruby",
	"domain_name", "doorkeeper", "dotenv", "dotenv-rails", "dry-configurable",
	"dry-container", "dry-core", "dry-inflector", "dry-logic", "dry-types",
	"elasticsearch", "elasticsearch-api", "elasticsearch-transport", "emoji_regex",
	"encryptor", "equalizer", "erubi", "erubis", "et-orbi", "ethon",
	"eventmachine", "excon", "execjs", "eye", "factory_bot", "factory_bot_rails",
	"faker", "faraday", "faraday-cookie_jar", "faraday-em_http",
	"faraday-em_synchrony", "faraday-excon", "faraday-http-cache",
	"faraday-httpclient", "faraday-multipart", "faraday-net_http",
	"faraday-net_http_persistent", "faraday-patron", "faraday-rack",
	"faraday-retry", "faraday_middleware", "faraday_middleware-aws-sigv4",
	"fastimage", "fastlane", "ffi", "ffi-compiler", "fluent-plugin-elasticsearch",
	"fluent-plugin-kubernetes_metadata_filter", "fluent-plugin-record-modifier",
	"fluent-plugin-s3", "fog-aws", "fog-core", "fog-google", "fog-json",
	"fog-local", "fog-xml", "foreman", "formatador", "fugit", "geocoder",
	"get_process_mem", "gh_inspector", "git", "gli", "globalid",
	"google-api-client", "google-apis-androidpublisher_v3", "google-apis-core",
	"google-apis-iamcredentials_v1", "google-apis-storage_v1", "google-cloud-core",
	"google-cloud-env", "google-cloud-errors", "google-cloud-storage",
	"google-protobuf", "googleapis-common-protos",
	"googleapis-common-protos-types", "googleauth", "graphql", "grpc", "guard",
	"guard-compat", "guard-rspec", "gyoku", "haml", "hashdiff", "hashie",
	"highline", "hike", "htmlentities", "http", "http-accept", "http-cookie",
	"http-form_data", "http_parser.rb", "httparty", "httpclient", "httpi", "i18n",
	"i18n_data", "ice_nine", "ipaddress", "jaro_winkler", "jbuilder", "jmespath",
	"jquery-rails", "jquery-ui-rails", "json", "json-jwt", "json-schema",
	"jsonapi-renderer", "jwt", "kaminari", "kaminari-actionview",
	"kaminari-activerecord", "kaminari-core", "kgio", "knapsack", "kostya-sigar",
	"kramdown", "kramdown-parser-gfm", "language_server-protocol", "launchy",
	"letter_opener", "libv8", "liquid", "listen", "little-plugger", "logging",
	"lograge", "logstash-filter-translate", "logstash-output-sqs", "loofah",
	"lumberjack", "mail", "marcel", "matrix", "memoist", "memory_profiler",
	"method_source", "mime-types", "mime-types-data", "mimemagic", "mini_magick",
	"mini_mime", "mini_portile2", "minitest", "molinillo", "money", "mongo",
	"msgpack", "multi_json", "multi_xml", "multipart-post", "mustermann", "mysql2",
	"nanaimo", "nap", "naturally", "nenv", "net-http-persistent", "net-imap",
	"net-ntp", "net-pop", "net-protocol", "net-scp", "net-sftp", "net-smtp",
	"net-ssh", "netrc", "newrelic_rpm", "nio4r", "no_proxy_fix", "nokogiri",
	"nori", "notiffany", "oauth", "oauth2", "octokit", "oj", "omniauth",
	"omniauth-google-oauth2", "omniauth-oauth2", "open4", "optimist",
	"orm_adapter", "os", "paper_trail", "paperclip", "parallel", "parallel_tests",
	"parser", "pdf-reader", "pg", "plist", "polyglot", "powerpack", "premailer",
	"premailer-rails", "pry", "pry-byebug", "pry-rails", "psych", "public_suffix",
	"puma", "pundit", "raabro", "racc", "rack", "rack-attack", "rack-cors",
	"rack-protection", "rack-proxy", "rack-test", "rack-timeout", "rails",
	"rails-controller-testing", "rails-deprecated_sanitizer", "rails-dom-testing",
	"rails-html-sanitizer", "rails-i18n", "railties", "rainbow", "raindrops",
	"rake", "ransack", "rb-fsevent", "rb-inotify", "rchardet", "rdoc", "redcarpet",
	"redis", "redis-actionpack", "redis-activesupport", "redis-namespace",
	"redis-rack", "redis-rails", "redis-store", "regexp_parser", "reline",
	"representable", "request_store", "responders", "rest-client", "retriable",
	"rexml", "roo", "rotp", "rouge", "rqrcode", "rspec", "rspec-core",
	"rspec-expectations", "rspec-its", "rspec-mocks", "rspec-rails", "rspec-retry",
	"rspec-support", "rspec_junit_formatter", "rubocop", "rubocop-ast",
	"rubocop-performance", "rubocop-rails", "rubocop-rspec", "ruby-macho",
	"ruby-progressbar", "ruby-rc4", "ruby-saml", "ruby2_keywords", "ruby_parser",
	"rubygems-bundler", "rubygems-update", "rubyzip", "rufus-scheduler",
	"safe_yaml", "sanitize", "sass", "sass-listen", "sass-rails", "sassc",
	"sassc-rails", "savon", "sawyer", "sdoc", "security", "selenium-webdriver",
	"sentry-raven", "sexp_processor", "shellany", "shoulda-matchers", "sidekiq",
	"signet", "simctl", "simple_form", "simplecov", "simplecov-html",
	"simplecov_json_formatter", "sinatra", "sixarm_ruby_unaccent",
	"slack-notifier", "slop", "spring", "spring-commands-rspec", "sprockets",
	"sprockets-rails", "sqlite3", "stackprof", "state_machines", "statsd-ruby",
	"stripe", "systemu", "temple", "term-ansicolor", "terminal-notifier",
	"terminal-table", "thin", "thor", "thread_safe", "tilt", "timecop", "timeout",
	"timers", "tins", "tomlrb", "trailblazer-option", "treetop", "ttfunk",
	"tty-cursor", "tty-screen", "tty-spinner", "turbolinks", "turbolinks-source",
	"twilio-ruby", "typhoeus", "tzinfo", "tzinfo-data", "uber", "uglifier", "unf",
	"unf_ext", "unicode-display_width", "unicode_utils", "unicorn",
	"uniform_notifier", "unparser", "uuidtools", "vcr", "virtus", "warden",
	"wasabi", "web-console", "webmock", "webpacker", "webrick", "websocket",
	"websocket-driver", "websocket-extensions", "will_paginate", "word_wrap",
	"xcodeproj", "xcpretty", "xcpretty-travis-formatter", "xml-simple", "xpath",
	"yajl-ruby", "yard", "zeitwerk",
	// Add more popular gems here
}

// generate variants of each popular gem
// complain if the name we received is one of the variants
func CheckForTyposquats(gemName string) []string {
	var matches []string

	for _, popular := range popularGems {
		// Skip if it's the exact same gem
		if strings.EqualFold(gemName, popular) {
			continue
		}

		variants := generateTyposquatVariants(popular)
		for _, variant := range variants {
			if strings.EqualFold(gemName, variant) {
				matches = append(matches, popular)
				break
			}
		}
	}

	return matches
}

// generateTyposquatVariants generates possible typosquat variations of a gem name
func generateTyposquatVariants(name string) []string {
	variants := make(map[string]bool)

	// Common typosquatting techniques:

	// 1. Character substitution
	substitutions := map[rune][]rune{
		'i': {'1', 'l'},
		'l': {'1', 'i'},
		'o': {'0'},
		'p': {'q'},
		'q': {'p'},
		'a': {'@', '4'},
		's': {'5', '$'},
		'e': {'3'},
		't': {'7'},
		'b': {'8'},
		'g': {'9'},
		'z': {'2'},
		'-': {'_'},
		'_': {'-'},
	}

	// Generate substitution variants
	runes := []rune(name)
	for i, r := range runes {
		if subs, ok := substitutions[r]; ok {
			for _, sub := range subs {
				newName := make([]rune, len(runes))
				copy(newName, runes)
				newName[i] = sub
				variants[string(newName)] = true
			}
		}
	}

	// 2. Missing character
	for i := range name {
		variant := name[:i] + name[i+1:]
		variants[variant] = true
	}

	// 3. Double character
	for i := range name {
		variant := name[:i] + string(name[i]) + name[i:]
		variants[variant] = true
	}

	// 4. Adjacent character swaps
	for i := 0; i < len(name)-1; i++ {
		runes := []rune(name)
		runes[i], runes[i+1] = runes[i+1], runes[i]
		variants[string(runes)] = true
	}

	// Convert map to slice
	result := make([]string, 0, len(variants))
	for variant := range variants {
		result = append(result, variant)
	}

	return result
}
