
class InvoiceItem {
  final String description;
  final int quantity;
  final double unitPrice;
  final double totalPrice;

  InvoiceItem({
    required this.description,
    required this.quantity,
    required this.unitPrice,
    required this.totalPrice,
  });

  factory InvoiceItem.fromJson(Map<String, dynamic> json) {
    return InvoiceItem(
      description: json['description'],
      quantity: json['quantity'],
      unitPrice: (json['unit_price'] as num).toDouble(),
      totalPrice: (json['total_price'] as num).toDouble(),
    );
  }
}

class Invoice {
  final String id;
  final String customerName;
  final String invoiceNumber;
  final DateTime issueDate;
  final double subtotal; 
  final double tax;      
  final double total;    
  final String status;
  final List<InvoiceItem>? items;

  Invoice({
    required this.id,
    required this.customerName,
    required this.invoiceNumber,
    required this.issueDate,
    required this.subtotal,
    required this.tax,
    required this.total,
    required this.status,
    this.items,
  });

  factory Invoice.fromJson(Map<String, dynamic> json) {
    return Invoice(
      id: json['id'],
      customerName: json['customer_name'] ?? 'N/A',
      invoiceNumber: json['invoice_number'],
      issueDate: DateTime.parse(json['issue_date']),
      subtotal: (json['subtotal'] as num).toDouble(),
      tax: (json['tax'] as num).toDouble(),
      total: (json['total'] as num).toDouble(),
      status: json['status'],
      items: json['items'] != null
          ? (json['items'] as List).map((item) => InvoiceItem.fromJson(item)).toList()
          : null,
    );
  }
}