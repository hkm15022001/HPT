import requests
import random
import string
provinces = [
    "Bắc Giang", "Bắc Kạn", "Cao Bằng", "Điện Biên", "Hà Giang", "Hà Nam", "Hải Dương", "Hải Phòng", "Hòa Bình", "Hưng Yên", "Lai Châu", "Lào Cai",
    "Nam Định", "Nghệ An", "Ninh Bình", "Phú Thọ", "Quảng Ninh", "Sơn La",
    "Thái Bình", "Thái Nguyên", "Thanh Hóa", "Tuyên Quang", "Vĩnh Phúc",
    "Yên Bái"
]
headers = {
    'Content-Type': 'application/json'
}

fixed_data = {
    "use_long_ship": True,
    "long_ship_id": 0,
    "long_ship_distance": 56,
    "transport_type_id": 3,
    "customer_send_id": 1,
    "customer_receive_id": 0,
    "detail": "Dien thoai iphone 12",
    "note": "Can than hang de vo",
    "sender": "Customer one - 231-233 Le Hong Phong - 6578678678",
    "receiver": "Customer two - 12 Bich Hoa - 6578678333"
}
def random_string(length):
    letters_and_digits = string.ascii_letters + string.digits
    return ''.join(random.choice(letters_and_digits) for i in range(length))
# Create 40 requests
for _ in range(40):
    random_province = random.choice(provinces)
    
    request_data = {
        "weight": random.randint(1, 30),  # Trọng lượng ngẫu nhiên từ 1 đến 30
        "volume": random.randint(1, 50),  # Thể tích ngẫu nhiên từ 1đến 50
        "type": random.choice(["Normal", "Express"]),  # Loại giao hàng ngẫu nhiên
        "receiver_address": random_province , # Tỉnh thành ngẫu nhiên từ danh sách
        "sender": random_string(20),  # Chuỗi ngẫu nhiên 20 kí tự và số
        "receiver": random_string(20)  # Chuỗi ngẫu nhiên 20 kí tự và số
    }

    request_data.update(fixed_data)

    response = requests.post('http://localhost:5000/api/order/create', headers=headers, json=request_data)
    
    if response.status_code == 201:
        print(f"Request { _ + 1 }: Success")
    else:
        print(f"Request { _ + 1 }: Failed with status code {response.status_code}")