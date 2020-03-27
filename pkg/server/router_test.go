package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupServer() *gin.Engine {
	s := &Server{
		Timeout:      1 * time.Second,
		MaxPerSecond: 5,
		SizeLimit:    10000,
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	gin.SetMode("test")

	return s.SetupRouter(r)
}

func TestPingRoute(t *testing.T) {

	router := SetupServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestFun(t *testing.T) {

	router := SetupServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `░░░█▀░░░░░░░░░░░▀▀███████░░░░
░░█▌░░░░░░░░░░░░░░░▀██████░░░
░█▌░░░░░░░░░░░░░░░░███████▌░░
░█░░░░░░░░░░░░░░░░░████████░░
▐▌░░░░░░░░░░░░░░░░░▀██████▌░░
░▌▄███▌░░░░▀████▄░░░░▀████▌░░
▐▀▀▄█▄░▌░░░▄██▄▄▄▀░░░░████▄▄░
▐░▀░░═▐░░░░░░══░░▀░░░░▐▀░▄▀▌▌
▐░░░░░▌░░░░░░░░░░░░░░░▀░▀░░▌▌
▐░░░▄▀░░░▀░▌░░░░░░░░░░░░▌█░▌▌
░▌░░▀▀▄▄▀▀▄▌▌░░░░░░░░░░▐░▀▐▐░
░▌░░▌░▄▄▄▄░░░▌░░░░░░░░▐░░▀▐░░
░█░▐▄██████▄░▐░░░░░░░░█▀▄▄▀░░
░▐░▌▌░░░░░░▀▀▄▐░░░░░░█▌░░░░░░
░░█░░▄▀▀▀▀▄░▄═╝▄░░░▄▀░▌░░░░░░
░░░▌▐░░░░░░▌░▀▀░░▄▀░░▐░░░░░░░
░░░▀▄░░░░░░░░░▄▀▀░░░░█░░░░░░░
░░░▄█▄▄▄▄▄▄▄▀▀░░░░░░░▌▌░░░░░░
░░▄▀▌▀▌░░░░░░░░░░░░░▄▀▀▄░░░░░
▄▀░░▌░▀▄░░░░░░░░░░▄▀░░▌░▀▄░░░
░░░░▌█▄▄▀▄░░░░░░▄▀░░░░▌░░░▌▄▄
░░░▄▐██████▄▄░▄▀░░▄▄▄▄▌░░░░▄░
░░▄▌████████▄▄▄███████▌░░░░░▄
░▄▀░██████████████████▌▀▄░░░░
▀░░░█████▀▀░░░▀███████░░░▀▄░░
░░░░▐█▀░░░▐░░░░░▀████▌░░░░▀▄░
░░░░░░▌░░░▐░░░░▐░░▀▀█░░░░░░░▀
░░░░░░▐░░░░▌░░░▐░░░░░▌░░░░░░░
░╔╗║░╔═╗░═╦═░░░░░╔╗░░╔═╗░╦═╗░
░║║║░║░║░░║░░░░░░╠╩╗░╠═╣░║░║░
░║╚╝░╚═╝░░║░░░░░░╚═╝░║░║░╩═╝░
`, w.Body.String())
}

func TestCompileRateLimitOk(t *testing.T) {

	router := SetupServer()

	helloworld := `println("Hello, world!")`

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(helloworld))
		req.Header.Set("X-Real-IP", "2601:7:1c82:4097:59a0:a80b:2841:b8c8")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestCompileRateLimitFail(t *testing.T) {

	router := SetupServer()

	helloworld := `println("Hello, world!")`

	countratelimited := 0

	for i := 0; i < 7; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(helloworld))
		req.Header.Set("X-Real-IP", "2601:7:1c82:4097:59a0:a80b:2841:b8c8")
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			countratelimited++
		}
	}
	assert.Equal(t, 2, countratelimited)
}

func TestTooLarge(t *testing.T) {

	router := SetupServer()

	toolarge := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris finibus velit quis nisi luctus mattis. Vivamus vitae ipsum scelerisque, auctor orci in, consectetur purus. Nam id orci tempor lectus posuere lobortis. Praesent dictum, lacus vitae vehicula suscipit, turpis nisl fringilla arcu, sed lobortis nunc tellus sed risus. Etiam consectetur lobortis erat, elementum tempus sem fermentum nec. Nulla porttitor, mauris in vehicula ullamcorper, nulla ligula ultricies sem, eget vehicula tortor lectus sit amet tortor. Curabitur id nisl in nibh iaculis placerat sed quis augue. Fusce cursus vitae leo et egestas. Praesent vitae nulla felis. Proin placerat massa neque, id aliquet velit feugiat et.
	Donec gravida mollis risus, quis luctus enim maximus at. Vivamus posuere gravida quam. Donec eget accumsan ante. Etiam pellentesque enim at tellus cursus, et sollicitudin dolor cursus. Cras risus felis, sagittis at ligula id, egestas euismod tellus. Curabitur fermentum, leo vel hendrerit sodales, mauris justo luctus nunc, eu ullamcorper risus tortor at magna. Sed efficitur rutrum sapien eu maximus. Nunc suscipit neque sapien, at fermentum ante lobortis sit amet. Mauris vitae fermentum justo, vitae posuere mauris. Nulla condimentum justo et viverra elementum.
	Aenean interdum nunc risus, a venenatis tellus suscipit sed. Vivamus ut porta diam, vitae commodo leo. Praesent ligula tellus, dictum et auctor in, hendrerit et tortor. Duis sem velit, malesuada id sem eget, tincidunt molestie dui. Curabitur eleifend sem tellus, ut porta lorem accumsan vitae. Pellentesque venenatis metus vitae felis iaculis, sed placerat mauris consequat. Nullam ut massa eu dui volutpat malesuada quis a sem.
	Maecenas at arcu nisi. Donec molestie nunc lacus, ac pharetra magna ultricies ac. Nunc tristique lectus sed velit consectetur varius. Nunc quis nulla et nunc tempor gravida. Nulla eget lacus et sem egestas facilisis id nec mauris. Pellentesque tincidunt, velit id iaculis pulvinar, purus enim molestie erat, non malesuada neque est a sapien. Maecenas facilisis sem ac accumsan elementum. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Nulla scelerisque congue nibh, quis congue turpis. Duis eu ante leo. Aliquam sed volutpat lectus.
	Proin fringilla neque eu arcu tincidunt molestie. Sed nulla ipsum, dapibus sed nisl sed, rutrum mattis magna. Vivamus aliquam risus at risus consequat, quis sodales est egestas. Sed blandit convallis iaculis. Praesent purus leo, dignissim sit amet est volutpat, consectetur luctus enim. Quisque aliquet non lacus quis pretium. Maecenas iaculis suscipit magna ac mattis. Morbi egestas, enim eget consectetur porttitor, risus ligula hendrerit metus, at ornare ante ex vel dui. Proin suscipit risus nec erat bibendum varius. Suspendisse at massa semper, cursus dui eget, varius dolor. Etiam a risus sed magna viverra interdum. Quisque congue lorem nibh.
	Sed congue laoreet turpis vel venenatis. Nulla imperdiet in velit et convallis. Pellentesque elit tortor, placerat nec efficitur eu, suscipit sed ligula. Fusce rhoncus semper tincidunt. Fusce mattis arcu nec porttitor aliquam. Duis aliquet laoreet urna sit amet commodo. Quisque eu dui at nisi egestas consectetur in id orci. Etiam non sodales odio, ut sollicitudin mi. Nullam finibus, sapien non placerat rhoncus, nibh mauris auctor urna, sed efficitur urna elit id nisi. Fusce bibendum sapien eu mattis euismod. Quisque sodales purus sed iaculis scelerisque. Sed in purus tincidunt orci faucibus sagittis.
	Interdum et malesuada fames ac ante ipsum primis in faucibus. Integer ut magna mauris. Proin a lorem purus. Maecenas quis orci elit. Sed fermentum tincidunt lectus, eu dictum ante molestie non. Mauris eget magna et tortor tempus sagittis. Duis efficitur erat vel tellus viverra faucibus eu aliquet risus. Aliquam sed blandit nisl, ac egestas ipsum. Nunc lacinia porta volutpat. Nam at tortor risus. Duis tempus metus vel purus vehicula tincidunt. Ut nunc metus, tincidunt in iaculis eget, cursus sed tellus. Nam dictum lorem et justo dignissim, sit amet aliquam felis cursus.
	Praesent ornare ligula ullamcorper nisi pellentesque congue. Nam egestas at augue at blandit. Aliquam tempor dui posuere, ultricies augue sed, ornare nibh. Integer feugiat purus nibh, et pretium felis mollis ac. Proin sit amet ultricies justo. Donec tempor ultricies tortor. Nullam odio dui, gravida eu accumsan vitae, ultrices et nulla.
	Quisque aliquet ipsum id arcu consectetur convallis. Sed ornare justo non semper lobortis. Sed eget purus condimentum, malesuada lectus sed, pharetra nulla. Integer blandit arcu sit amet ante malesuada, a efficitur tortor consequat. Fusce lobortis eros sit amet magna porta, eget aliquam justo ultrices. Sed imperdiet molestie tincidunt. Pellentesque pulvinar tincidunt ex, at pulvinar sem interdum nec. Vivamus interdum mattis pharetra. Integer porttitor lacus eu est semper, vitae congue purus ultricies. Fusce eu molestie lorem. Curabitur cursus interdum magna. Nullam nunc enim, vehicula sed viverra ut, volutpat vel massa.
	Suspendisse tempus sapien id dolor iaculis laoreet. Curabitur semper nulla non cursus tincidunt. Cras pulvinar lobortis lacus, id suscipit nisi tincidunt eu. Aliquam a leo est. Nam consectetur id sapien at faucibus. Proin faucibus justo ut massa elementum, a varius eros tincidunt. Sed posuere elit vel ipsum consectetur, sed aliquam ante tristique. Etiam vitae fermentum ante, sed lobortis magna. Etiam sit amet gravida nisi. Fusce finibus imperdiet diam, eu maximus justo. Integer porta lectus vel mauris consectetur viverra. Mauris pretium diam tellus, in vulputate diam imperdiet ut. Phasellus suscipit pellentesque bibendum.
	Nulla sit amet ligula felis. Aliquam congue lectus nulla, a venenatis odio efficitur ut. Nam sagittis libero tristique felis placerat placerat. Nulla consequat mauris ac tellus luctus volutpat. Nullam sed sapien neque. Integer a lacus et turpis accumsan volutpat. Donec nec condimentum odio, commodo viverra ante. Maecenas at ex ultrices, bibendum mauris malesuada, posuere justo. Aenean suscipit vehicula eros, quis efficitur dui fringilla sit amet. Donec sed odio tempus, aliquam velit id, consequat ligula. Nulla facilisi. Nunc elementum semper dui, sed placerat velit sagittis nec. Nunc sed euismod neque. Donec et arcu eget ante condimentum blandit sed id risus.
	Donec quis ligula vulputate, convallis libero vitae, sollicitudin turpis. Aliquam sed condimentum turpis, nec maximus risus. Pellentesque sodales aliquet leo. Sed ac tincidunt arcu. Fusce venenatis id ligula imperdiet sagittis. Aliquam a arcu tortor. Nullam sed sapien sed diam ullamcorper dapibus.
	Nunc id neque eu odio posuere ultrices. In hac habitasse platea dictumst. In dolor nulla, efficitur nec auctor eu, tristique non turpis. Etiam eget lorem ultricies, sagittis neque eget, consectetur ex. Donec eleifend porta pellentesque. Suspendisse ut orci in erat aliquam laoreet. Quisque laoreet metus ac mi pretium mollis. Aenean convallis tellus ut massa facilisis, in volutpat metus sagittis.
	Suspendisse vel turpis ut ex rutrum sodales nec non velit. In volutpat interdum odio, eget posuere velit aliquam ac. Etiam tristique diam a nisl accumsan, in aliquam metus maximus. Mauris nec accumsan lorem. Nam at auctor velit. Phasellus neque dui, varius non nisi eget, porta tempor velit. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Pellentesque dolor augue, finibus id vehicula in, efficitur et diam. Duis pretium sed orci vel lobortis. Phasellus mauris tortor, feugiat maximus tristique volutpat, pellentesque eu nulla. Morbi bibendum nisi sed risus placerat aliquam. Ut sed volutpat dolor, et malesuada elit. Aenean laoreet, orci ut fringilla tincidunt, tortor risus auctor elit, vel efficitur leo tellus quis massa.
	Curabitur venenatis dictum lobortis. Morbi ac elementum augue, at rutrum risus. Sed ultrices sem non lectus mollis, eget facilisis mi mollis. Nullam lacinia eget risus sit amet accumsan. Nullam luctus eros ligula, at convallis elit venenatis in. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. In hac habitasse platea dictumst.
	Quisque viverra volutpat orci, ut sodales velit dapibus eu. Quisque sed metus at dolor auctor gravida. In ac convallis est, at pretium orci. Aliquam erat volutpat. Pellentesque sollicitudin, ante nec vulputate convallis, velit arcu mattis erat, nec blandit odio arcu id urna. Cras id condimentum odio, a gravida lacus. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Morbi vehicula vel tortor in fermentum. Nam sit amet sem sed purus suscipit interdum. Nullam vel nisl nulla. Sed ornare tortor ut dui sollicitudin vehicula.
	Curabitur molestie eleifend eros at pulvinar. Mauris varius ipsum eu lorem aliquet tincidunt. Donec eget tincidunt enim, et porta nulla. Donec consectetur tempus risus sit amet sollicitudin. In hac habitasse platea dictumst. Nam vel massa in justo rhoncus auctor ut porta metus. Phasellus imperdiet arcu ex, ut rhoncus sem dictum ut. Cras ultricies eu nulla in porta.
	Quisque elementum nibh at turpis sollicitudin varius. Duis fringilla dictum nibh, at semper orci rutrum eu. Vestibulum molestie mollis sagittis. In magna est, varius vitae justo quis, volutpat porta nunc. Etiam fermentum, erat ut imperdiet sollicitudin, ex lorem porttitor arcu, eget imperdiet justo metus id arcu. Mauris eleifend pharetra diam, sed rhoncus metus finibus a. Aliquam erat volutpat. Etiam volutpat quam non luctus efficitur. In eu maximus arcu. Aenean ut urna sed sapien tincidunt condimentum. Etiam nec nisi condimentum augue euismod vehicula ut scelerisque erat. In elementum volutpat egestas. Mauris elit lacus, tempus eget imperdiet eu, feugiat at ipsum.
	Cras gravida risus sit amet tempus faucibus. In semper dui ut volutpat dapibus. Nulla lobortis, diam sed pellentesque fermentum, diam justo suscipit tortor, sit amet fringilla dolor augue tincidunt augue volutpat.`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(toolarge))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestHelloWorld(t *testing.T) {

	router := SetupServer()

	helloworld := `println("Hello, world!")`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(helloworld))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Hello, world!\n", response.Output)
	assert.Equal(t, 0, response.Status)
	assert.Empty(t, response.Errors)
}

func TestError(t *testing.T) {
	router := SetupServer()

	wtf := "wtf"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(wtf))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Empty(t, response.Output)
	assert.NotEqual(t, 0, response.Status)
	assert.NotEmpty(t, response.Errors)
}

func TestTimeout(t *testing.T) {
	router := SetupServer()

	wtf := `
	while true {
		println("bitch")
	}
	`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(wtf))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Empty(t, response.Output)
	assert.Equal(t, -1, response.Status)
	assert.NotEmpty(t, response.Errors)
}

func TestFibonacci(t *testing.T) {
	router := SetupServer()

	input := `
	func fib(n int) int {
		if n == 0 {
		return 0
		}

		if n == 1 {
		return 1
		}

		return fib(n-1) + fib(n-2)
	}

	println(fib(30))`

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(input))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 0, response.Status)
	assert.Equal(t, "832040\n", response.Output)
	assert.Empty(t, response.Errors)
}
