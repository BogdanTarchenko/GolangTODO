import SwiftUI

struct CreateTaskView: View {
    @Environment(\.dismiss) private var dismiss
    @ObservedObject var viewModel: TasksViewModel
    
    @State private var title: String = ""
    @State private var description: String = ""
    @State private var deadline: Date = Date()
    @State private var priority: TaskPriority? = nil
    @State private var showDatePicker: Bool = false
    @State private var showDateError: Bool = false
    @State private var showTitleError: Bool = false
    @State private var showMacroError: Bool = false
    
    var body: some View {
        NavigationView {
            Form {
                Section("Основная информация") {
                    TextField("Название", text: $title)
                        .onChange(of: title) { _ in
                            showTitleError = false
                            showMacroError = false
                        }
                    if showTitleError {
                        Text("Название должно содержать минимум 4 символа")
                            .foregroundColor(.red)
                            .font(.caption)
                    }
                    if showMacroError {
                        Text("Некорректный формат даты в макросе. Используйте формат DD.MM.YYYY")
                            .foregroundColor(.red)
                            .font(.caption)
                    }
                    TextField("Описание", text: $description)
                }
                
                Section("Приоритет") {
                    Toggle("Установить приоритет", isOn: Binding(
                        get: { priority != nil },
                        set: { isOn in
                            if isOn {
                                priority = .medium
                            } else {
                                priority = nil
                            }
                        }
                    ))
                    
                    if priority != nil {
                        Picker("Приоритет", selection: Binding(
                            get: { priority ?? .medium },
                            set: { priority = $0 }
                        )) {
                            Text("Низкий").tag(TaskPriority.low)
                            Text("Средний").tag(TaskPriority.medium)
                            Text("Высокий").tag(TaskPriority.high)
                            Text("Критический").tag(TaskPriority.critical)
                        }
                        .pickerStyle(.menu)
                    }
                }
                
                Section("Дедлайн") {
                    Toggle("Установить дедлайн", isOn: $showDatePicker)
                    if showDatePicker {
                        DatePicker(
                            "Дедлайн",
                            selection: $deadline,
                            in: Date()...,
                            displayedComponents: [.date, .hourAndMinute]
                        )
                    }
                }
            }
            .navigationTitle("Новая задача")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Отмена") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Создать") {
                        createTask()
                    }
                    .disabled(title.isEmpty)
                }
            }
            .alert("Ошибка", isPresented: $showDateError) {
                Button("OK", role: .cancel) { }
            } message: {
                Text("Дедлайн не может быть в прошлом")
            }
        }
    }
    
    private func validateMacros() -> Bool {
        let pattern = #"!before (\d{2}\.\d{2}\.\d{4})"#
        guard let regex = try? NSRegularExpression(pattern: pattern) else { return true }
        
        let range = NSRange(title.startIndex..., in: title)
        let matches = regex.matches(in: title, range: range)
        
        for match in matches {
            guard let dateRange = Range(match.range(at: 1), in: title) else { continue }
            let dateString = String(title[dateRange])
            
            let formatter = DateFormatter()
            formatter.dateFormat = "dd.MM.yyyy"
            formatter.locale = Locale(identifier: "ru_RU")
            
            guard let date = formatter.date(from: dateString) else {
                showMacroError = true
                return false
            }
            
            if date <= Date() {
                showMacroError = true
                return false
            }
        }
        
        return true
    }
    
    private func createTask() {
        if title.count < 4 {
            showTitleError = true
            return
        }
        
        if !validateMacros() {
            return
        }
        
        let deadlineString = showDatePicker ? formatDate(deadline) : nil
        
        if showDatePicker && deadline < Date() {
            showDateError = true
            return
        }
        
        viewModel.createTask(
            title: title,
            description: description.isEmpty ? nil : description,
            deadline: deadlineString,
            priority: priority
        )
        
        dismiss()
    }
    
    private func formatDate(_ date: Date) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime]
        return formatter.string(from: date)
    }
}
