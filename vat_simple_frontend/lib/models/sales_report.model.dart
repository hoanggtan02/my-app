class SalesDataPoint {
  final DateTime date;
  final double revenue;

  SalesDataPoint({required this.date, required this.revenue});

  factory SalesDataPoint.fromJson(Map<String, dynamic> json) {
    return SalesDataPoint(
      date: DateTime.parse(json['date']),
      revenue: (json['revenue'] as num).toDouble(),
    );
  }
}