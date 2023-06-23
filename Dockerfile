
FROM scratch
ENTRYPOINT [ "/geo2tz" ]
CMD [ "start" ]
COPY /geo2tz /geo2tz
COPY /tzdata /tzdata