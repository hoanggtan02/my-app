import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:intl/intl.dart';
import 'package:vat_simple_frontend/models/sales_report.model.dart';
import 'package:vat_simple_frontend/utils/constants.dart';

class ReportService {
  Future<List<SalesDataPoint>> getSalesReport(String token, DateTime startDate, DateTime endDate) async {
    final dateFormat = DateFormat('yyyy-MM-dd');
    final url = Uri.parse(
        '$apiUrl/reports/sales?start_date=${dateFormat.format(startDate)}&end_date=${dateFormat.format(endDate)}');
        
    try {
      final response = await http.get(
        url,
        headers: {'Authorization': 'Bearer $token'},
      );

      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(utf8.decode(response.bodyBytes));
        return data.map((item) => SalesDataPoint.fromJson(item)).toList();
      } else {
        throw Exception('Failed to load sales report');
      }
    } catch (e) {
      throw Exception('Could not connect to server.');
    }
  }
}