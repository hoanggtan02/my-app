import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:vat_simple_frontend/models/user.model.dart'; 
import 'package:vat_simple_frontend/utils/constants.dart';

class AuthService {
  // Hàm đăng nhập
  Future<String> loginUser(String email, String password) async {
    final url = Uri.parse('$apiUrl/auth/login');

    try {
      final response = await http.post(
        url,
        headers: {
          'Content-Type': 'application/json; charset=UTF-8',
        },
        body: jsonEncode({
          'email': email,
          'password': password,
        }),
      );

      final responseData = jsonDecode(response.body);

      if (response.statusCode == 200) {
        final String token = responseData['access_token'];
        return token;
      } else {
        throw Exception(responseData['message'] ?? 'Đã có lỗi xảy ra.');
      }
    } catch (error) {
      print(error);
      throw Exception('Không thể kết nối đến server.');
    }
  }

  // Hàm lấy thông tin user
  Future<User> getUserProfile(String token) async {
    final url = Uri.parse('$apiUrl/users/me');

    try {
      final response = await http.get(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token', // Gửi token để xác thực
        },
      );

      final responseData = jsonDecode(response.body);

      if (response.statusCode == 200) {
        return User.fromJson(responseData);
      } else {
        throw Exception(responseData['message'] ?? 'Failed to load user profile.');
      }
    } catch (error) {
      print(error);
      throw Exception('Could not connect to the server.');
    }
  }

  Future<void> registerUser(
      String email, String password, String companyName) async {
    final url = Uri.parse('$apiUrl/auth/register');
    try {
      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json; charset=UTF-8'},
        body: jsonEncode({
          'email': email,
          'password': password,
          'company_name': companyName,
        }),
      );

      if (response.statusCode == 201) {
        // Đăng ký thành công, không cần làm gì thêm
        return;
      } else {
        // Nếu có lỗi, ví dụ email đã tồn tại
        final responseData = jsonDecode(response.body);
        throw Exception(responseData['message'] ?? 'Đăng ký thất bại.');
      }
    } catch (error) {
      throw Exception('Không thể kết nối đến server.');
    }
  }
}