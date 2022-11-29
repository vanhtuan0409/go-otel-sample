.PHONY: template

template:
	@helm template --release-name local \
		--namespace tracing \
		.

diff:
	make -s template | kubectl diff --namespace tracing -f -

apply:
	make -s template | kubectl apply --namespace tracing -f -

dryrun:
	make -s template | kubectl apply --namespace tracing -f - --dry-run=server
