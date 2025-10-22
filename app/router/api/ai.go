package api

import (
	"net/http"
	"strings"
	"time"

	"student_shared/app/model"
	"student_shared/app/utils/ai"
	"student_shared/app/utils/database"

	"github.com/gin-gonic/gin"
	// 引入统一的请求包
	req "student_shared/app/model/req"
	"gorm.io/gorm"
)


// SummarizeText 生成摘要与关键词，并在有 note_id 时持久化到 NoteAIMeta
func SummarizeText(c *gin.Context) {
	// 需要认证
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	userID := userIDVal.(uint)
	roleVal, _ := c.Get("role")
	role := ""
	if roleVal != nil { role = roleVal.(string) }

	var params req.SummarizeRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	var text string
	var noteID uint

	if params.NoteID != nil {
		// 加载笔记并进行权限校验
		note, err := model.GetNoteByID(database.DB, *params.NoteID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
			return
		}
		// 私有笔记仅作者或管理员可处理
		if note.Status == "private" && !(note.UserID == userID || role == "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权处理该私有笔记"})
			return
		}
		text = strings.TrimSpace(note.Content)
		noteID = note.ID
	} else if params.Content != nil {
		text = strings.TrimSpace(*params.Content)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "至少提供 note_id 或 content"})
		return
	}

	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容为空，无法生成摘要"})
		return
	}

	summary, keywords := ai.Summarize(text)

	// 如有 note_id，写入/更新 NoteAIMeta
	if noteID != 0 {
		kw := strings.Join(keywords, ",")
		now := time.Now()
		if err := model.UpsertNoteAIMeta(database.DB, noteID, summary, kw, now); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存AI元数据失败"})
			return
		}
	}

	// 返回摘要与关键词
	c.JSON(http.StatusOK, gin.H{"summary": summary, "keywords": keywords})
}

// 获取笔记的AI元数据（摘要/关键词），公开笔记可匿名查看，私有笔记仅作者或管理员
func GetNoteAIMeta(c *gin.Context) {
	var params req.GetNoteRequest
	if err := c.ShouldBindJSON(&params); err != nil || params.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}
	note, err := model.GetNoteByID(database.DB, params.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}
	role := ""
	if rv, ok := c.Get("role"); ok { role = rv.(string) }
	userID := uint(0)
	if uv, ok := c.Get("userID"); ok { userID = uv.(uint) }
	if note.Status == "private" && !(note.UserID == userID || role == "admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权查看该私有笔记的AI信息"})
		return
	}
	meta, err := model.GetNoteAIMetaByNoteID(database.DB, note.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"note_id": note.ID, "summary": "", "keywords": []string{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询AI元数据失败"})
		return
	}
	keywords := []string{}
	if strings.TrimSpace(meta.Keywords) != "" {
		keywords = strings.Split(meta.Keywords, ",")
	}
	c.JSON(http.StatusOK, gin.H{
		"note_id": note.ID,
		"summary": meta.Summary,
		"keywords": keywords,
		"plagiarism_score": meta.PlagiarismScore,
		"flags": meta.Flags,
		"last_reviewed_at": meta.LastReviewedAt,
	})
}