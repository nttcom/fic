# Copyright 2020 NTT Limited and NTT Communications Corporation All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

fmt:
	gofmt -s -w .

fmtcheck:
	(! gofmt -s -d . | grep '^')

test: fmtcheck
	go test ./... -count=1

build: test
	gox -os="linux darwin windows" -arch="386 amd64" -output="./build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}"

archive:
	find ./build -name "*_linux_*"  -type f -and ! -name "*.tar.gz" -exec tar -zcvf {}.tar.gz {} \;
	find ./build ! -name "*_linux_*"  -type f -and ! -name "*.zip" -exec zip {}.zip {} \;

clean:
	rm -f ./build/*

doc:
	rm -f ./doc/*
	go run doc.go

.PHONY: fmt fmtcheck test build archive clean doc
