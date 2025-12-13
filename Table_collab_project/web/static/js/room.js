class TableCollabRoom {
	constructor() {
		this.roomId = window.location.pathname.split('/').pop()
		this.userId = null
		this.ws = null
		this.participants = new Map()

		this.init()
	}

	init() {
		document.getElementById('roomId').textContent = this.roomId

		this.username =
			localStorage.getItem('username') ||
			'User_' + Math.random().toString(36).substr(2, 4)

		this.connectWebSocket()
		this.setupEventListeners()
	}

	connectWebSocket() {
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
		const wsUrl = `${protocol}//${window.location.host}/ws/${this.roomId}`

		this.ws = new WebSocket(wsUrl)

		this.ws.onopen = () => {
			console.log('Connected to room:', this.roomId)
			this.sendJoin()
		}

		this.ws.onmessage = event => {
			this.handleMessage(JSON.parse(event.data))
		}

		this.ws.onclose = () => {
			console.log('Disconnected')
		}

		this.ws.onerror = error => {
			console.error('WebSocket error:', error)
		}
	}

	sendJoin() {
		this.ws.send(
			JSON.stringify({
				type: 'join_room',
				payload: {
					username: this.username,
					color: this.getUserColor(),
				},
			})
		)
	}

	handleMessage(data) {
		console.log('Received:', data)

		switch (data.type) {
			case 'join_room':
				this.addParticipant(data.user_id, data.payload)
				break

			case 'leave_room':
				this.removeParticipant(data.user_id)
				break

			case 'cursor_move':
				this.updateCursor(data.user_id, data.payload)
				break

			case 'text_update':
				this.updateText(data.payload)
				break

			case 'chat_message':
				this.addChatMessage(data.user_id, data.payload)
				break
		}
	}

	addParticipant(userId, data) {
		this.participants.set(userId, data)
		this.updateParticipantsList()
	}

	removeParticipant(userId) {
		this.participants.delete(userId)
		this.updateParticipantsList()
	}

	updateParticipantsList() {
		const list = document.getElementById('participants')
		const count = document.getElementById('participantCount')

		list.innerHTML = ''
		this.participants.forEach((data, userId) => {
			const div = document.createElement('div')
			div.className = 'participant'
			div.innerHTML = `
                <div style="color: ${data.color}">●</div>
                <div>${data.username || 'Anonymous'}</div>
            `
			list.appendChild(div)
		})

		count.textContent = `${this.participants.size} participants`
	}

	getUserColor() {
		const colors = ['#FF6B6B', '#4ECDC4', '#FFD166', '#06D6A0']
		return colors[Math.floor(Math.random() * colors.length)]
	}

	updateCursor(userId, position) {
		// Implementation for cursor tracking
	}

	updateText(payload) {
		const editor = document.getElementById('editor')
		if (editor.value !== payload.text) {
			editor.value = payload.text
		}
	}

	addChatMessage(userId, payload) {
		const chat = document.getElementById('chatMessages')
		const message = document.createElement('div')
		message.className = 'chat-message'
		message.innerHTML = `<strong>${
			this.participants.get(userId)?.username || 'Unknown'
		}:</strong> ${payload.text}`
		chat.appendChild(message)
		chat.scrollTop = chat.scrollHeight
	}

	setupEventListeners() {
		const editor = document.getElementById('editor')
		const chatInput = document.getElementById('chatInput')
		const sendBtn = document.getElementById('sendBtn')

		let debounceTimer
		editor.addEventListener('input', e => {
			clearTimeout(debounceTimer)
			debounceTimer = setTimeout(() => {
				this.sendTextUpdate(e.target.value)
			}, 300)
		})

		document.addEventListener('mousemove', e => {
			this.sendCursorMove(e.clientX, e.clientY)
		})

		const sendMessage = () => {
			const text = chatInput.value.trim()
			if (text && this.ws) {
				this.ws.send(
					JSON.stringify({
						type: 'chat_message',
						payload: { text: text },
					})
				)
				chatInput.value = ''
			}
		}

		sendBtn.addEventListener('click', sendMessage)
		chatInput.addEventListener('keypress', e => {
			if (e.key === 'Enter') sendMessage()
		})
	}

	sendTextUpdate(text) {
		if (this.ws) {
			this.ws.send(
				JSON.stringify({
					type: 'text_update',
					payload: { text: text, version: 1 },
				})
			)
		}
	}

	sendCursorMove(x, y) {
		if (this.ws) {
			this.ws.send(
				JSON.stringify({
					type: 'cursor_move',
					payload: { x: x, y: y },
				})
			)
		}
	}
}

// Initialize when page loads
window.addEventListener('DOMContentLoaded', () => {
	new TableCollabRoom()
})
