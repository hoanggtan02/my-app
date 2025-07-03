class CreateInvoiceItemRequest {
  final String productId;
  final int quantity;

  CreateInvoiceItemRequest({required this.productId, required this.quantity});

  Map<String, dynamic> toJson() => {
        'product_id': productId,
        'quantity': quantity,
      };
}

class CreateInvoiceRequest {
  final String customerId;
  final String invoiceNumber;
  final DateTime issueDate;
  final DateTime dueDate;
  final List<CreateInvoiceItemRequest> items;

  CreateInvoiceRequest({
    required this.customerId,
    required this.invoiceNumber,
    required this.issueDate,
    required this.dueDate,
    required this.items,
  });

  Map<String, dynamic> toJson() => {
        'customer_id': customerId,
        'invoice_number': invoiceNumber,
        'issue_date': issueDate.toIso8601String(),
        'due_date': dueDate.toIso8601String(),
        'items': items.map((item) => item.toJson()).toList(),
      };
}