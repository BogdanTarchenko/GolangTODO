//
//  TaskDTO.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

import Foundation

struct CreateTaskDTO: Codable {
    var title: String
    var description: String?
    var deadline: Date?
    var priority: TaskPriority?
}

struct UpdateTaskDTO: Codable {
    var title: String?
    var description: String?
    var deadline: Date?
    var status: TaskStatus?
    var priority: TaskPriority?
}
