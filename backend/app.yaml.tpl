api:
  server_address: ":9090"
  context_timeout: 2
database: 
  host: ${POSTGRES_HOST}
  port: ${POSTGRES_PORT}
  user: ${POSTGRES_USER}
  password: ${POSTGRES_PASSWORD}
  name: ${POSTGRES_DB}
jwt:
  private_key: |
    -----BEGIN PRIVATE KEY-----
    MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDcxhNwFwJtXSdI
    QE+6Ioe82HyYTZpn6KadsNI5lzcdpgvj9ngPFUbqGtSJarMc93qznwk2tuZkKC8k
    cOOkrJMBwZ9bLK8W+euRMHgk9vZEuGB+/v6n8N3mvW+vsn6g0fg7LKa2gILhs1uK
    eiPXfneqL7bCk3xjU7/Q8CDxlE9OLEwkICphD37ll2y4PRcrchzhO579+M6YrBMI
    G9DTPWt2DvueZQyXZdWudio9UmzgIqHLjoeyROlSIO1KX10U3CGv2UTT6YsUoeyp
    8wjozLSKXR7RlzZ53NXrCXdLwE6V6NTeQEqXdKb4yPS5oDyqPFTi39Hfx7heGfBo
    lf3FNJ4bAgMBAAECggEAYoyep6X1xOjUvKlMjYuVaOSANaJKfwC4w2JnbSLFjRwO
    abufGyiFx8GjRxYUjyUfpiejPsPFM0dGx+8Ghv8r/hg2sMXRAKIeF+j5cJK3GrTt
    CjN8bG4WN8YvIVA9uz8PHicP4h6ajfJ4tedQsYR4GUWEQPYCC/qaAMP4CK6J+hvy
    v2vMMbzN1PHBKXIeax5qWHpXMvOVGMOWdsZ8Rc7UbgY80dKYNs+ahgLOxiVHxsXs
    aoFJnIe81p24RLHtabH1N8cGhZOyINjnPvoeNE/HW7LZFJ/J5k9kuIrGtOAAKZjh
    bit2qtnR5fVO0pVPsdrxYVTL/2M5Hd3oLAf4mFeF3QKBgQDsOtUWDKeDxma4Twkq
    h43xJtiWfIFatQVELwho5NodEfh9ZfL06wlxbsJ8t2jJpcQKQ3zDWej3MKi9W4Qu
    FVnaUfj1sL50EVH2Syxh5HPRrXX4R/hni0jdJVjCahHHE2u1dYnoziogDbPdMWc/
    xzEqAj6/ocSBjb7dHoM/DXG+FQKBgQDvQBnJAufcKer/10xmylf+AZGKyPHqSw2+
    9SzPIY6o8OiTG6x/0ldHOLGnZeAiivQq4v9A76YFvPhz9UKSv00VXxKr86Uu+gKO
    c4VbVzwe7/Fed++m/PP3uQKo8JY/+/nBbfJCvFlL0l9rLo+TEQYAoXAsX/a0b2EZ
    lOrgErYnbwKBgH7Ef48KkWZstLjZaQDSp4AuqXHwNHZZyA6z8p5fmRCakS+x4vQ9
    oN6nYmUNA4WamB4t4yjt+c+U5ChhkQgt2v8GmEQ4aavdk49I/fM2ZlSx8imfbZUb
    MKnEHeKOiyW6rUU+Yxh0cjSrRcdAeLjICwERHV020T34s+DzO9k9PLmVAoGBAKbZ
    BSJxrFCVyxTwiI+GvSae4WjwCgViogtx3/XzaRHYL9mnivz5K3S3zOz41v4/+VeP
    RoN6nUWTK5FykSLV1mP5EYRpPeEs6Wt+lJnGlF7e5m0DJ1ZFQb6Yf4phfebRSrPi
    gPiZcYy3AWQ17FqbnJwD+b54jgv3QLgeak4pvm5xAoGAUmiDI0Jbqi+5UdxgiLxT
    pOXxK31rp4OBsLXCF2pMJteWGF4nqRjhawB5si8Qp6AVVlfVK05CFmuIAIDZzUMb
    Kd3u0fmLDiDKxWHieyfKirJ5lF0FD194zaY0ndn1gR2AztbUkEWDLK6heVA39AJS
    K1YJybLpwnAmqgy1hTfLLMg=
    -----END PRIVATE KEY-----
  public_key: |
    -----BEGIN PUBLIC KEY-----
    MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3MYTcBcCbV0nSEBPuiKH
    vNh8mE2aZ+imnbDSOZc3HaYL4/Z4DxVG6hrUiWqzHPd6s58JNrbmZCgvJHDjpKyT
    AcGfWyyvFvnrkTB4JPb2RLhgfv7+p/Dd5r1vr7J+oNH4OyymtoCC4bNbinoj1353
    qi+2wpN8Y1O/0PAg8ZRPTixMJCAqYQ9+5ZdsuD0XK3Ic4Tue/fjOmKwTCBvQ0z1r
    dg77nmUMl2XVrnYqPVJs4CKhy46HskTpUiDtSl9dFNwhr9lE0+mLFKHsqfMI6My0
    il0e0Zc2edzV6wl3S8BOlejU3kBKl3Sm+Mj0uaA8qjxU4t/R38e4XhnwaJX9xTSe
    GwIDAQAB
    -----END PUBLIC KEY-----
  token_expiry: 1440m
  permissions_claim: "permissions"
  serial_claim: "serial"

