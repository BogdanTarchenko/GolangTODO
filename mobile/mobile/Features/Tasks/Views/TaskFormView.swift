//
//  TaskFormView.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

import SwiftUI

struct TaskFormView: View {
    @Environment(\.dismiss) private var dismiss
    let viewModel: TaskListViewModel
    
    @State private var title = ""
    @State private var description = ""
    @State private var deadline: Date?
    @State private var priority: TaskPriority = .medium
    @State private var showingDatePicker = false
    
    var body: some View {
        NavigationView {
            Form {
                Section("Основное") {
                    TextField("Название", text: $title)
                        .autocapitalization(.sentences)
                    TextField("Описание", text: $description)
                        .autocapitalization(.sentences)
                }
                
                Section("Дополнительно") {
                    HStack {
                        Text("Дедлайн")
                        Spacer()
                        Button(action: { showingDatePicker.toggle() }) {
                            Text(deadline?.formatted(date: .abbreviated, time: .omitted) ?? "Не задан")
                                .foregroundColor(deadline == nil ? .secondary : .primary)
                        }
                    }
                    
                    if showingDatePicker {
                        DatePicker("", selection: Binding(
                            get: { deadline ?? Date() },
                            set: { deadline = $0 }
                        ), displayedComponents: .date)
                        .datePickerStyle(.graphical)
                        
                        Button("Убрать дедлайн") {
                            deadline = nil
                            showingDatePicker = false
                        }
                        .foregroundColor(.red)
                    }
                    
                    Picker("Приоритет", selection: $priority) {
                        Text("Критический").tag(TaskPriority.critical)
                        Text("Высокий").tag(TaskPriority.high)
                        Text("Средний").tag(TaskPriority.medium)
                        Text("Низкий").tag(TaskPriority.low)
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
                        viewModel.createTask(
                            title: title,
                            description: description.isEmpty ? nil : description,
                            deadline: deadline,
                            priority: priority
                        )
                        dismiss()
                    }
                    .disabled(title.count < 4)
                }
            }
        }
    }
}

#Preview {
    TaskFormView(viewModel: TaskListViewModel())
}
