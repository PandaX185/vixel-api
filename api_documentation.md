# Vixel API Documentation

## Register User
### Request
POST /api/v1/users
Body: {"username":"testuser","email":"test@example.com","password":"password123"}
### Response
```json
{"data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NzAyNzI1NTMsInN1YiI6MX0.2Y3tphyR4nj0hWyrnCUEAbie2gusqhnhbPV7HHprZ2w"},"status":"resource created","timestamp":"2026-02-05T07:22:33.129842217+02:00"}
```

## Login
### Request
POST /api/v1/users/login
Body: {"email":"test@example.com","password":"password123"}
### Response
```json
{"data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NzAyNzI1NTMsInN1YiI6MX0.2Y3tphyR4nj0hWyrnCUEAbie2gusqhnhbPV7HHprZ2w"},"status":"success","timestamp":"2026-02-05T07:22:33.215603657+02:00"}
```

## Upload Image
### Request
POST /api/v1/images
Headers: Authorization: Bearer {token}
Form: file=@test.jpg, alt_text=Test Image
### Response
```json
{"data":{"id":1,"url":"http://localhost:9000/vixel/vixel-948407-1770268953248950612","alt_text":"Test Image","user_id":1},"status":"resource created","timestamp":"2026-02-05T07:22:33.255874787+02:00"}
```

## Get Image
### Request
GET /api/v1/images/{id}
Headers: Authorization: Bearer {token}
### Response
```json
{"data":{"id":2,"url":"http://localhost:9000/vixel/vixel-64229-1770268953268247785","alt_text":"Get Test","user_id":1},"status":"success","timestamp":"2026-02-05T07:22:33.289651263+02:00"}
```

## List User Images
### Request
GET /api/v1/users/{user_id}/images
Headers: Authorization: Bearer {token}
### Response
```json
{"data":[{"id":1,"url":"http://localhost:9000/vixel/vixel-948407-1770268953248950612","alt_text":"Test Image","user_id":1},{"id":2,"url":"http://localhost:9000/vixel/vixel-64229-1770268953268247785","alt_text":"Get Test","user_id":1}],"status":"success","timestamp":"2026-02-05T07:22:33.298415631+02:00"}
```


# Image Transformations

## Transform Image - Resize
### Request
POST /api/v1/images/{id}/transform
Headers: Authorization: Bearer {token}, Content-Type: application/json
Body: {"resize":{"width":50,"height":50}}
### Response
```json
{"data":{"new_image_url":"http://localhost:9000/vixel/vixel-531156-1770268953335339626"},"status":"success","timestamp":"2026-02-05T07:22:33.340098355+02:00"}
```

## Transform Image - Crop
### Request
POST /api/v1/images/{id}/transform
Headers: Authorization: Bearer {token}, Content-Type: application/json
Body: {"crop":{"x":10,"y":10,"width":50,"height":50}}
### Response
```json
{"data":{"new_image_url":"http://localhost:9000/vixel/vixel-840938-1770268953374432133"},"status":"success","timestamp":"2026-02-05T07:22:33.379468937+02:00"}
```

## Transform Image - Rotate
### Request
POST /api/v1/images/{id}/transform
Headers: Authorization: Bearer {token}, Content-Type: application/json
Body: {"rotate":{"angle":90}}
### Response
```json
{"data":{"new_image_url":"http://localhost:9000/vixel/vixel-177372-1770268953416129915"},"status":"success","timestamp":"2026-02-05T07:22:33.421564019+02:00"}
```

## Transform Image - Flip Horizontal
### Request
POST /api/v1/images/{id}/transform
Headers: Authorization: Bearer {token}, Content-Type: application/json
Body: {"flip":{"direction":"horizontal"}}
### Response
```json
{"data":{"new_image_url":"http://localhost:9000/vixel/vixel-984432-1770268953456714267"},"status":"success","timestamp":"2026-02-05T07:22:33.462148371+02:00"}
```

## Transform Image - Flip Vertical
### Request
POST /api/v1/images/{id}/transform
Headers: Authorization: Bearer {token}, Content-Type: application/json
Body: {"flip":{"direction":"vertical"}}
### Response
```json
{"data":{"new_image_url":"http://localhost:9000/vixel/vixel-160837-1770268953498266155"},"status":"success","timestamp":"2026-02-05T07:22:33.50334598+02:00"}
```

## Transform Image - Format Conversion
### Request
POST /api/v1/images/{id}/transform
Headers: Authorization: Bearer {token}, Content-Type: application/json
Body: {"format_conversion":{"format":"png"}}
### Response
```json
{"data":{"new_image_url":"http://localhost:9000/vixel/vixel-55917-1770268953539526051"},"status":"success","timestamp":"2026-02-05T07:22:33.544246218+02:00"}
```


## Transform Image - Filter
### Request
POST /api/v1/images/{id}/transform
Headers: Authorization: Bearer {token}, Content-Type: application/json
Body: {"filter":{"saturation":50,"brightness":10,"contrast":20}}
### Response
```json
{"data":{"new_image_url":"http://localhost:9000/vixel/vixel-697113-1770268953579569803"},"status":"success","timestamp":"2026-02-05T07:22:33.584548357+02:00"}
```


## Transform Image - Watermark
### Request
POST /api/v1/images/{id}/transform
Headers: Authorization: Bearer {token}, Content-Type: application/json
Body: {"watermark":{"text":"Test","position":{"x":10,"y":10},"opacity":50}}
### Response
```json

```


## Delete Image
### Request
DELETE /api/v1/images/{id}
Headers: Authorization: Bearer {token}
### Response
```json
{"data":{"message":"image deleted"},"status":"success","timestamp":"2026-02-05T07:13:24.180781809+02:00"}
```


