import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:vat_simple_frontend/models/customer.model.dart';
import 'package:vat_simple_frontend/utils/constants.dart';

class CustomerService {
  Future<List<Customer>> getCustomers(String token) async {
    final url = Uri.parse('$apiUrl/customers');
    try {
      final response = await http.get(
        url,
        headers: {
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final List<dynamic> responseData = jsonDecode(response.body);
        return responseData.map((data) => Customer.fromJson(data)).toList();
      } else {
        throw Exception('Failed to load customers');
      }
    } catch (e) {
      throw Exception('Could not connect to server.');
    }
  }
}