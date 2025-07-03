import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:vat_simple_frontend/models/invoice.model.dart';
import 'package:vat_simple_frontend/utils/constants.dart';
import 'package:vat_simple_frontend/models/order_item.model.dart';


class InvoiceService {
  // Lấy danh sách hóa đơn
  Future<List<Invoice>> getInvoices(String token) async {
    final url = Uri.parse('$apiUrl/invoices');
    try {
      final response = await http.get(
        url,
        headers: {'Authorization': 'Bearer $token'},
      );
      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(utf8.decode(response.bodyBytes));
        return data.map((item) => Invoice.fromJson(item)).toList();
      } else {
        throw Exception('Failed to load invoices');
      }
    } catch (e) {
      throw Exception('Could not connect to server');
    }
  }

  // Tạo hóa đơn mới
  Future<void> createInvoice(String token, List<OrderItem> items) async {
    final url = Uri.parse('$apiUrl/invoices');
    
    final List<Map<String, dynamic>> itemsJson = items.map((item) => {
      'product_id': item.product.id,
      'quantity': item.quantity
    }).toList();

    try {
      final response = await http.post(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
        body: jsonEncode({'items': itemsJson}),
      );
      if (response.statusCode != 201) {
        throw Exception('Failed to create invoice');
      }
    } catch (e) {
      throw Exception('Could not connect to server.');
    }
  }
   Future<Invoice> getInvoiceDetails(String token, String invoiceId) async {
    final url = Uri.parse('$apiUrl/invoices/$invoiceId');
    try {
      final response = await http.get(
        url,
        headers: {'Authorization': 'Bearer $token'},
      );
      if (response.statusCode == 200) {
        // Dùng utf8.decode để đảm bảo đọc đúng tiếng Việt
        return Invoice.fromJson(json.decode(utf8.decode(response.bodyBytes)));
      } else {
        throw Exception('Failed to load invoice details');
      }
    } catch (e) {
      throw Exception('Could not connect to server.');
    }
  }
}