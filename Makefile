all:
	go build github.com/geonetwork/gnm_collect
	mkdir -p gnm_collect_1.0
	if [ -f gnm_collect ]; then\
		cp gnm_collect gnm_collect_1.0/;\
	fi

	if [ -f gnm_collect.exe ]; then \
		cp gnm_collect.exe gnm_collect_1.0/ ;\
	fi

	cp -r ../../gonum/plot/vg/fonts gnm_collect_1.0/
	tar czf gnm_collect_1.0.tar.gz gnm_collect_1.0/
	rm -rf gnm_collect_1.0
