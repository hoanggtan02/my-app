import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:vat_simple_frontend/models/customer.model.dart';
import 'package:vat_simple_frontend/utils/constants.dart';

class CustomerService {
  // Lấy danh sách khách hàng
  Future<List<Customer>> getCustomers(String token) async {
    final url = Uri.parse('$apiUrl/customers');
    try {
      final response = await http.get(
        url,
        headers: {'Authorization': 'Bearer $token'},
      );
      if (response.statusCode == 200) {
        final List<dynamic> responseData = jsonDecode(utf8.decode(response.bodyBytes));
        return responseData.map((data) => Customer.fromJson(data)).toList();
      } else {
        throw Exception('Failed to load customers');
      }
    } catch (e) {
      throw Exception('Could not connect to server.');
    }
  }

  // Tạo khách hàng mới
  Future<void> createCustomer(String token, Map<String, String> customerData) async {
    final url = Uri.parse('$apiUrl/customers');
    try {
      final response = await http.post(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
        body: jsonEncode(customerData),
      );
      if (response.statusCode != 201) {
        throw Exception('Failed to create customer');
      }
    } catch (e) {
      throw Exception('Could not connect to server.');
    }
  }

  // Cập nhật khách hàng
  Future<void> updateCustomer(String token, String customerId, Map<String, String> customerData) async {
    final url = Uri.parse('$apiUrl/customers/$customerId');
    try {
      final response = await http.patch(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
        body: jsonEncode(customerData),
      );
      if (response.statusCode != 200) {
        throw Exception('Failed to update customer');
      }
    } catch (e) {
      throw Exception('Could not connect to server.');
    }
  }

  // Xóa khách hàng
  Future<void> deleteCustomer(String token, String customerId) async {
    final url = Uri.parse('$apiUrl/customers/$customerId');
    try {
      final response = await http.delete(
        url,
        headers: {'Authorization': 'Bearer $token'},
      );
      if (response.statusCode != 204) {
        throw Exception('Failed to delete customer');
      }
    } catch (e) {
      throw Exception('Could not connect to server.');
    }
  }
}