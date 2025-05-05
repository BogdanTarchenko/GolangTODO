//
//  Task.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

import Foundation

enum TaskStatus: String, Codable {
    case active = "ACTIVE"
    case completed = "COMPLETED"
    case overdue = "OVERDUE"
    case late = "LATE"
}

enum TaskPriority: String, Codable {
    case low = "LOW"
    case medium = "MEDIUM"
    case high = "HIGH"
    case critical = "CRITICAL"
}

struct Task: Identifiable, Codable {
    let id: Int
    var title: String
    var description: String?
    var deadline: Date?
    var status: TaskStatus
    var priority: TaskPriority
    let createdAt: Date
    var updatedAt: Date?
    
    enum CodingKeys: String, CodingKey {
        case id
        case title
        case description
        case deadline
        case status
        case priority
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}
