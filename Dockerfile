FROM gcr.io/distroless/static-debian12:nonroot
COPY uri-template-tester /
COPY public /public/
USER nonroot:nonroot
ENTRYPOINT ["/uri-template-tester"]
EXPOSE 8080
