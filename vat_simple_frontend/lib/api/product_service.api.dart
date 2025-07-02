import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:vat_simple_frontend/models/product.model.dart';
import 'package:vat_simple_frontend/utils/constants.dart';

class ProductService {
  // Lấy danh sách sản phẩm
  Future<List<Product>> getProducts(String token) async {
    final url = Uri.parse('$apiUrl/products');
    try {
      final response = await http.get(
        url,
        headers: {'Authorization': 'Bearer $token'},
      );
      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(utf8.decode(response.bodyBytes));
        return data.map((item) => Product.fromJson(item)).toList();
      } else {
        throw Exception('Failed to load products');
      }
    } catch (e) {
      throw Exception('Could not connect to server');
    }
  }

  // Tạo sản phẩm mới
  Future<void> createProduct(String token, String name, String? description, double unitPrice) async {
    final url = Uri.parse('$apiUrl/products');
    try {
      final response = await http.post(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
        body: jsonEncode({
          'name': name,
          'description': description,
          'unit_price': unitPrice,
        }),
      );
      if (response.statusCode != 201) {
        throw Exception('Failed to create product');
      }
    } catch (e) {
      throw Exception('Could not connect to server');
    }
  }
}