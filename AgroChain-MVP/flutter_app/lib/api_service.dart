import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  // Get the server root URL from the environment variable
  static const String _serverUrl = String.fromEnvironment(
    'BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  // --- THIS LINE IS CORRECTED ---
  // Append the /api prefix to the server URL
  static const String _baseUrl = "$_serverUrl/api";
  // --- END OF CORRECTION ---

  Future<Map<String, dynamic>> createProduce(String crop, int quantity, String owner) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/produce'),
      headers: <String, String>{'Content-Type': 'application/json; charset=UTF-8'},
      body: jsonEncode(<String, dynamic>{'crop': crop, 'quantity': quantity, 'owner': owner}),
    );
    if (response.statusCode == 201) {
      return jsonDecode(response.body);
    } else {
      throw Exception('Failed to create produce: ${response.statusCode} ${response.reasonPhrase}');
    }
  }

  Future<List<dynamic>> getProduceHistory(String id) async {
    final response = await http.get(Uri.parse('$_baseUrl/produce/$id/history'));
    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    } else {
      throw Exception('Failed to load history: ${response.statusCode} ${response.reasonPhrase}');
    }
  }
}