# Phoenix Game Client: Main Controller
# This script acts as the "Eyes" of the game, visualizing the Python simulation state.
extends Control

var go_stdout: FileAccess
var agent_sprite: Sprite2D
var goal_sprite: Sprite2D
var complete_label: Label
var code_editor: CodeEdit
var error_console: RichTextLabel
var level_title: Label
var description_label: RichTextLabel
var obstacles_container: Node2D
var telemetry_label: Label
var trail_container: Node2D
var trail_positions = []
const MAX_TRAIL = 5

var current_lang = "py"
var current_script_path = ""
var active_pipe = null

var templates = {
	"py": "# Phoenix Agent Loop\nwhile get_x() < get_goal_x():\n    move_forward()\n",
	"cpp": "#include \"agent.h\"\nvoid agent_main() {\n    while(get_x() < get_goal_x()) {\n        move_forward();\n    }\n}\n",
	"java": "public class Agent {\n    public static void main(String[] args) {\n        // Agent logic\n    }\n}"
}

func _ready():
	agent_sprite = $HBoxContainer/WorldView/SimulationSpace/Agent
	goal_sprite = $HBoxContainer/WorldView/SimulationSpace/Goal
	obstacles_container = $HBoxContainer/WorldView/SimulationSpace/Obstacles
	trail_container = $HBoxContainer/WorldView/SimulationSpace/Trail
	complete_label = $HBoxContainer/WorldView/UIMessage
	code_editor = $HBoxContainer/EditorPanel/CodeEditor
	error_console = $HBoxContainer/EditorPanel/ErrorConsole
	level_title = $HBoxContainer/WorldView/LevelTitle
	description_label = $HBoxContainer/EditorPanel/DescriptionLabel
	telemetry_label = $HBoxContainer/WorldView/TelemetryPanel/TelemetryLabel
	
	setup_highlighter()
	complete_label.hide()
	set_language("python")
	start_python_core()

func setup_highlighter():
	var highlighter = CodeHighlighter.new()
	
	# Default text color (bright white/light grey for dark background)
	code_editor.add_theme_color_override("font_color", Color(0.9, 0.9, 0.9, 1))
	
	# Cyan keywords
	var kw_color = Color(0, 1, 1, 1)
	var keywords = ["def", "class", "if", "elif", "else", "while", "for", "in", "return", "import", "void", "int", "public", "static"]
	for kw in keywords:
		highlighter.add_keyword_color(kw, kw_color)
		
	# Bright Green API Functions
	var api_color = Color(0, 1, 0.2, 1)
	var api_funcs = ["move_forward", "move_backward", "move_up", "move_down", "get_x", "get_y", "get_goal_x", "get_goal_y", "moveForward", "moveBackward", "moveUp", "moveDown", "getX", "getY", "getGoalX", "getGoalY"]
	for fn in api_funcs:
		highlighter.add_keyword_color(fn, api_color)
		
	# Light Grey comments (brighter than the background)
	highlighter.add_color_region("#", "", Color(0.6, 0.6, 0.6, 1), true)
	highlighter.add_color_region("//", "", Color(0.6, 0.6, 0.6, 1), true)
	
	# Amber strings
	highlighter.add_color_region('"', '"', Color(1, 0.75, 0, 1), false)
	
	# Numbers and Symbols (Pink/Purple)
	highlighter.number_color = Color(0.8, 0.4, 0.8, 1)
	highlighter.symbol_color = Color(0.9, 0.9, 0.9, 1)
	highlighter.function_color = Color(0.4, 0.8, 1, 1)
	highlighter.member_variable_color = Color(0.9, 0.9, 0.9, 1)
	
	code_editor.syntax_highlighter = highlighter

func set_language(lang):
	current_lang = lang
	code_editor.text = templates[lang]
	# Update button visual state
	var toggle_box = $HBoxContainer/EditorPanel/LangToggle
	for child in toggle_box.get_children():
		if child is Button:
			if child.name == "BtnPython" and lang == "python": child.modulate = Color(0, 1, 1, 1)
			elif child.name == "BtnCPP" and lang == "cpp": child.modulate = Color(0, 1, 1, 1)
			elif child.name == "BtnJava" and lang == "java": child.modulate = Color(0, 1, 1, 1)
			else: child.modulate = Color(0.5, 0.5, 0.5, 1)

func _on_lang_python(): set_language("python")
func _on_lang_cpp(): set_language("cpp")
func _on_lang_java(): set_language("java")

func start_python_core():
	var python_bin = "python3"
	var core_path = ProjectSettings.globalize_path("res://../core/main.py")
	
	var ext = "py"
	if current_lang == "cpp": ext = "cpp"
	if current_lang == "java": ext = "java"
	
	current_script_path = ProjectSettings.globalize_path("res://../core/scripts/agent." + ext)

	# Get the level path from our global store, default to level_1 if missing
	var level_path = ProjectSettings.get_setting("game/current_level_path", ProjectSettings.globalize_path("res://../core/levels/level_1.json"))

	OS.set_environment("PHX_SCRIPT_PATH", current_script_path)
	OS.set_environment("PHX_LEVEL_PATH", level_path)

	active_pipe = OS.execute_with_pipe(python_bin, [core_path], true)
	
	if active_pipe.has("stdio"):
		go_stdout = active_pipe["stdio"]
		print("Pipe opened successfully to Python core.")
	else:
		print("Failed to open pipe to Python core.")

func _process(_delta):
	if go_stdout and go_stdout.get_length() > 0:
		var line = go_stdout.get_line().strip_edges()
		if line.begins_with("{"):
			var json = JSON.new()
			var error = json.parse(line)
			if error == OK:
				var data = json.get_data()
				update_ui(data)
			else:
				print("JSON Parse Error: ", json.get_error_message(), " in line: ", line)

func update_ui(data):
	# 1. Handle Level Metadata (Tick 0 only)
	if data.has("level"):
		var level = data["level"]
		level_title.text = level["title"]
		description_label.text = "[center]" + level["description"] + "[/center]"
		print("Level Loaded: ", level["title"])
		
		# Update API Reference
		var api_list = $HBoxContainer/EditorPanel/ApiReferencePanel/ApiList
		for child in api_list.get_children():
			child.queue_free()
		
		if level.has("allowed_builtins"):
			for api_fn in level["allowed_builtins"]:
				var lbl = Label.new()
				lbl.text = api_fn + "() → void"
				if api_fn.begins_with("get"):
					lbl.text = api_fn + "() → int"
				lbl.add_theme_color_override("font_color", Color(0, 1, 0.2, 1))
				api_list.add_child(lbl)
		
		# Draw obstacles
		for child in obstacles_container.get_children():
			child.queue_free()
			
		if level.has("obstacles"):
			for obs in level["obstacles"]:
				var rect = ColorRect.new()
				rect.color = Color(0.8, 0.2, 0.2, 0.8) # Red walls
				rect.position = Vector2(obs["x"] / 100.0, obs["y"] / 100.0)
				rect.size = Vector2(obs["w"] / 100.0, obs["h"] / 100.0)
				obstacles_container.add_child(rect)

	# 2. Handle Goal Position
	if data.has("goal"):
		var goal_pos = data["goal"]
		goal_sprite.position = Vector2(goal_pos["x"] / 100.0, goal_pos["y"] / 100.0)

	# 3. Handle Agent Position
	if data.has("agent"):
		var pos = data["agent"]
		var screen_pos = Vector2(pos["x"] / 100.0, pos["y"] / 100.0)
		
		# Add to trail before updating agent
		var dot = ColorRect.new()
		dot.size = Vector2(8, 8)
		dot.position = agent_sprite.position - Vector2(4, 4)
		dot.color = Color(0, 1, 0.62, 0.5) # Cyan trail
		trail_container.add_child(dot)
		trail_positions.append(dot)
		if trail_positions.size() > MAX_TRAIL:
			var old_dot = trail_positions.pop_front()
			old_dot.queue_free()
			
		# Fade trail
		for i in range(trail_positions.size()):
			trail_positions[i].color.a = float(i + 1) / MAX_TRAIL * 0.5
			
		agent_sprite.position = screen_pos
		
		# Update Telemetry
		if data.has("tick") and data.has("status") and data.has("goal"):
			var goal_p = data["goal"]
			telemetry_label.text = "TICK: %03d   STATUS: %s   X: %05d   Y: %05d   GOAL: %05d, %05d" % [data["tick"], data["status"].to_upper(), pos["x"], pos["y"], goal_p["x"], goal_p["y"]]
		
	# 4. Handle Status
	if data.has("status"):
		if data["status"] == "complete":
			complete_label.text = "MISSION COMPLETE\nSolved in " + str(data["tick"]) + " ticks"
			complete_label.add_theme_color_override("font_color", Color(0, 1, 0, 1))
			complete_label.show()
			error_console.hide()
		elif data["status"] == "crashed":
			complete_label.text = "CRASHED INTO WALL\nTick: " + str(data["tick"])
			complete_label.add_theme_color_override("font_color", Color(1, 0, 0, 1))
			complete_label.show()
			error_console.hide()
		elif data["status"] == "compile_error" or data["status"] == "runtime_error":
			complete_label.hide()
			error_console.text = "[color=red]" + data.get("error_msg", "Unknown error") + "[/color]"
			error_console.show()
		else:
			complete_label.hide()
			error_console.hide()

func kill_core():
	if active_pipe:
		if active_pipe.has("stdio") and is_instance_valid(active_pipe["stdio"]):
			active_pipe["stdio"].close()
		if active_pipe.has("pid"):
			var pid = active_pipe["pid"]
			# Send SIGTERM to the process
			OS.kill(pid)
		active_pipe = null

func _on_deploy_pressed():
	print("Deploying new script...")
	var new_code = code_editor.text
	
	var ext = "py"
	if current_lang == "cpp": ext = "cpp"
	if current_lang == "java": ext = "java"
	
	var target_script_path = ProjectSettings.globalize_path("res://../core/scripts/agent." + ext)
	var tmp_path = target_script_path + ".tmp"
	
	# Atomic write: save to .tmp first
	var file = FileAccess.open(tmp_path, FileAccess.WRITE)
	if file:
		file.store_string(new_code)
		file.close()
		
		# Move/Rename file over the target to prevent partial reads by the Python core
		var dir = DirAccess.open(tmp_path.get_base_dir())
		if dir:
			if dir.file_exists(target_script_path.get_file()):
				dir.remove(target_script_path.get_file())
			dir.rename(tmp_path.get_file(), target_script_path.get_file())
			print("Script atomically saved to disk: ", target_script_path)
	
	if target_script_path != current_script_path:
		print("Language changed. Restarting core...")
		kill_core()
		start_python_core()

func _on_back_pressed():
	print("Returning to Level Select...")
	kill_core()
	get_tree().change_scene_to_file("res://level_select.tscn")

func _on_reset_pressed():
	print("Resetting simulation...")
	kill_core()
	start_python_core()

func toggle_fullscreen():
	var current_mode = DisplayServer.window_get_mode()
	if current_mode == DisplayServer.WINDOW_MODE_FULLSCREEN:
		DisplayServer.window_set_mode(DisplayServer.WINDOW_MODE_WINDOWED)
	else:
		DisplayServer.window_set_mode(DisplayServer.WINDOW_MODE_FULLSCREEN)

func _unhandled_input(event):
	if event is InputEventKey and event.pressed and event.keycode == KEY_F11:
		toggle_fullscreen()
