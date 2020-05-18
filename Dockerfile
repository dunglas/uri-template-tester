FROM gcr.io/distroless/static
COPY uri-template-tester /
COPY public /public/
CMD ["/uri-template-tester"]
EXPOSE 80 443
