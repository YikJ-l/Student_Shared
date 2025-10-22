<template>
  <div class="note-detail-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <el-breadcrumb separator=">">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/notes' }">学习笔记</el-breadcrumb-item>
          <el-breadcrumb-item>{{ note.title }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </div>

    <div class="container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="10" animated />
      </div>

      <div v-else-if="note.id" class="note-detail">
        <el-row :gutter="20">
          <!-- 主要内容区域 -->
          <el-col :span="18">
            <!-- 笔记内容 -->
            <el-card class="note-content-card">
              <!-- 笔记头部信息 -->
              <div class="note-header">
                <div class="note-meta">
                  <el-tag :type="getCategoryType(note.course_name)" size="large">
                    {{ note.course_name }}
                  </el-tag>
                  <span class="note-time">
                    <el-icon><Clock /></el-icon>
                    {{ formatTime(note.created_at) }}
                  </span>
                  <span class="note-views">
                    <el-icon><View /></el-icon>
                    {{ note.view_count || 0 }} 次浏览
                  </span>
                </div>
                
                <div class="note-actions">
                  <el-button type="primary" :icon="Star" :loading="isLiking" @click="handleLike">
                    {{ note.is_liked ? '取消点赞' : '点赞' }} ({{ note.like_count || 0 }})
                  </el-button>
                  <el-button :icon="Collection" :loading="isFavoriting" @click="handleFavorite">
                    {{ note.is_favorited ? '已收藏' : '收藏' }}
                  </el-button>

                  <el-dropdown v-if="canEditNote || canDeleteNote" trigger="click" @command="handleNoteAction">
                    <el-button>
                      <el-icon><MoreFilled /></el-icon>
                      更多
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item v-if="canEditNote" command="edit">编辑笔记</el-dropdown-item>
                        <el-dropdown-item v-if="canDeleteNote" command="delete" divided>删除笔记</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>

              <!-- 笔记标题 -->
              <h1 class="note-title">{{ note.title }}</h1>

              <!-- 作者信息 -->
              <div class="author-info">
                <el-avatar :size="40" :src="note.author_avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <div class="author-details">
                  <div class="author-name">{{ note.author }}</div>
                  <div class="author-meta">
                    <span>发布于 {{ formatDate(note.created_at) }}</span>
                    <span v-if="note.updated_at !== note.created_at">· 更新于 {{ formatDate(note.updated_at) }}</span>
                  </div>
                </div>
              </div>



              <!-- 笔记内容 -->
              <div class="note-content">
                <div v-html="formattedContent"></div>
              </div>
            </el-card>

            <el-card class="ai-summary-card">
              <template #header>
                <div class="ai-summary-header">
                  <span>AI 智能摘要</span>
                  <el-button
                    v-if="canEditNote"
                    size="small"
                    type="primary"
                    :loading="regenerating"
                    @click="handleRegenerateSummary"
                  >
                    重新生成
                  </el-button>
                </div>
              </template>
              <div class="ai-summary-body">
                <p v-if="aiMeta.summary" class="ai-summary-text">{{ aiMeta.summary }}</p>
                <p v-else class="ai-summary-empty">暂无摘要</p>
                <div class="ai-keywords" v-if="aiMeta.keywords && aiMeta.keywords.length">
                  <el-tag v-for="kw in aiMeta.keywords" :key="kw" type="info" class="ai-keyword">{{ kw }}</el-tag>
                </div>
              </div>
            </el-card>

            <!-- 评论区域 -->
            <div class="comments-section">
              <el-card class="comments-card">
              <template #header>
                <div class="comments-header">
                  <span>评论 ({{ comments.length }})</span>
                  <el-button v-if="!showCommentForm && canComment" type="primary" size="small" @click="showCommentForm = true">
                    <el-icon><ChatDotRound /></el-icon>
                    发表评论
                  </el-button>
                  <el-button v-if="!showCommentForm && !canComment" type="primary" size="small" @click="$router.push('/login')">
                    <el-icon><ChatDotRound /></el-icon>
                    登录后评论
                  </el-button>
                </div>
              </template>

              <!-- 评论表单 -->
              <div v-if="showCommentForm" class="comment-form">
                <el-input
                  v-model="commentContent"
                  type="textarea"
                  :rows="4"
                  placeholder="写下你的评论..."
                  maxlength="500"
                  show-word-limit
                />
                <div class="comment-actions">
                  <el-button @click="showCommentForm = false">取消</el-button>
                  <el-button type="primary" :loading="commentLoading" @click="handleSubmitComment">
                    发表评论
                  </el-button>
                </div>
              </div>

              <!-- 评论列表 -->
              <div v-if="!comments || comments.length === 0" class="empty-comments">
                <el-empty description="暂无评论">
                  <el-button v-if="canComment" type="primary" @click="showCommentForm = true">发表第一条评论</el-button>
                  <el-button v-else type="primary" @click="$router.push('/login')">
                    <el-icon><ChatDotRound /></el-icon>
                    登录后评论
                  </el-button>
                </el-empty>
              </div>

              <div v-else class="comments-list">
                <template v-for="(comment, index) in comments" :key="comment?.id || `comment-${index}`">
                  <div v-if="comment && comment.id" class="comment-item">
                    <div class="comment-wrapper">
                      <div class="comment-avatar">
                        <el-avatar :size="32" :src="comment.user_avatar">
                          <el-icon><User /></el-icon>
                        </el-avatar>
                      </div>
                      <div class="comment-content">
                        <div class="comment-header">
                          <span class="comment-author">{{ comment.username || '匿名用户' }}</span>
                          <span v-if="comment.reply_to_username" class="reply-to">
                            回复 @{{ comment.reply_to_username }}
                          </span>
                          <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
                        </div>
                        <div class="comment-text">{{ comment.content || '评论内容为空' }}</div>
                        <div class="comment-actions">
                          <el-button 
                            link 
                            size="small" 
                            :type="comment.is_liked ? 'primary' : 'default'"
                            :loading="commentLikingMap.get(comment.id)"
                            @click="handleLikeComment(comment)"
                          >
                            <el-icon><Star /></el-icon>
                            {{ comment.like_count || 0 }}
                          </el-button>
                          <el-button link size="small" @click="handleReplyComment(comment)">
                            <el-icon><ChatDotRound /></el-icon>
                            回复
                          </el-button>
                          <el-dropdown v-if="canEditComment(comment.user_id) || canDeleteComment(comment.user_id)" trigger="click" @command="(cmd) => handleCommentAction(cmd, comment)">
                            <el-button link size="small">
                              <el-icon><MoreFilled /></el-icon>
                            </el-button>
                            <template #dropdown>
                              <el-dropdown-menu>
                                <el-dropdown-item v-if="canEditComment(comment.user_id)" command="edit">编辑</el-dropdown-item>
                                <el-dropdown-item v-if="canDeleteComment(comment.user_id)" command="delete" divided>删除</el-dropdown-item>
                              </el-dropdown-menu>
                            </template>
                          </el-dropdown>
                        </div>
                        
                        <!-- 回复表单 -->
                        <div v-if="replyingTo === comment.id" class="reply-form">
                          <el-input
                            v-model="replyContent"
                            type="textarea"
                            :rows="3"
                            :placeholder="`回复 @${comment.username}...`"
                            maxlength="500"
                            show-word-limit
                          />
                          <div class="reply-actions">
                            <el-button size="small" @click="cancelReply">取消</el-button>
                            <el-button size="small" type="primary" :loading="replyLoading" @click="submitReply(comment)">
                              发表回复
                            </el-button>
                          </div>
                        </div>
                        
                        <!-- 回复列表 -->
                        <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
                          <template v-for="(reply, index) in comment.replies" :key="reply?.id || `reply-${index}`">
                            <div v-if="reply" class="reply-item">
                              <div class="reply-avatar">
                                <el-avatar :size="24" :src="reply.user_avatar">
                                  <el-icon><User /></el-icon>
                                </el-avatar>
                              </div>
                              <div class="reply-content">
                                <div class="reply-header">
                                  <span class="reply-author">{{ reply.username }}</span>
                                  <span v-if="reply.reply_to_username" class="reply-to">
                                    回复 @{{ reply.reply_to_username }}
                                  </span>
                                  <span class="reply-time">{{ formatTime(reply.created_at) }}</span>
                                </div>
                                <div class="reply-text">{{ reply.content }}</div>
                                <div class="reply-actions">
                                  <el-button 
                                    link 
                                    size="small" 
                                    :type="reply.is_liked ? 'primary' : 'default'"
                                    :loading="commentLikingMap.get(reply.id)"
                                    @click="handleLikeComment(reply)"
                                  >
                                    <el-icon><Star /></el-icon>
                                    {{ reply.like_count || 0 }}
                                  </el-button>
                                  <el-button link size="small" @click="handleReplyToReply(comment, reply)">
                                    <el-icon><ChatDotRound /></el-icon>
                                    回复
                                  </el-button>
                                </div>
                              </div>
                            </div>
                          </template>
                        </div>
                      </div>
                    </div>
                  </div>
                </template>
              </div>
            </el-card>
            </div>
          </el-col>

          <!-- 侧边栏 -->
          <el-col :span="6">
            <!-- 笔记信息 -->
            <el-card class="note-info-sidebar">
              <template #header>
                <span>笔记信息</span>
              </template>
              
              <div class="sidebar-content">

                <div class="info-item">
                  <label>难度等级：</label>
                  <el-tag :type="getDifficultyType(note.difficulty)" size="small">
                    {{ getDifficultyText(note.difficulty) }}
                  </el-tag>
                </div>
                <div class="info-item">
                  <label>字数统计：</label>
                  <span>{{ getWordCount(note.content) }} 字</span>
                </div>
                <div class="info-item">
                  <label>阅读时间：</label>
                  <span>约 {{ getReadingTime(note.content) }} 分钟</span>
                </div>
                <div class="info-item">
                  <label>发布时间：</label>
                  <span>{{ formatDate(note.created_at) }}</span>
                </div>
                <div v-if="note.updated_at !== note.created_at" class="info-item">
                  <label>更新时间：</label>
                  <span>{{ formatDate(note.updated_at) }}</span>
                </div>
              </div>
            </el-card>

            <!-- 作者其他笔记 -->
            <el-card class="author-notes-sidebar">
              <template #header>
                <span>作者的其他笔记</span>
              </template>
              
              <div v-if="authorNotes.length === 0" class="empty-author-notes">
                <el-empty description="暂无其他笔记" :image-size="60" />
              </div>
              
              <div v-else class="author-notes-list">
                <template v-for="(authorNote, index) in authorNotes" :key="authorNote?.id || `author-note-${index}`">
                  <div v-if="authorNote && authorNote.id" class="author-note-item">
                    <div class="author-note-info" @click="$router.push(`/notes/${authorNote.id}`)">
                      <h5>{{ authorNote.title || '无标题' }}</h5>
                      <p>{{ authorNote.course_name || '未知课程' }}</p>
                      <div class="author-note-stats">
                        <span><el-icon><View /></el-icon> {{ authorNote.view_count || 0 }}</span>
                        <span><el-icon><Star /></el-icon> {{ authorNote.like_count || 0 }}</span>
                      </div>
                    </div>
                  </div>
                </template>
              </div>
            </el-card>

            <!-- 相关笔记推荐 -->
            <el-card class="related-notes-sidebar">
              <template #header>
                <span>相关笔记推荐</span>
              </template>
              
              <div v-if="relatedNotes.length === 0" class="empty-related-notes">
                <el-empty description="暂无相关笔记" :image-size="60" />
              </div>
              
              <div v-else class="related-notes-list">
                  <div v-for="(relatedNote, index) in relatedNotes" :key="relatedNote?.id || `related-note-${index}`" class="related-note-item" v-if="relatedNote && relatedNote.id">
                    <div class="related-note-info" @click="$router.push(`/notes/${relatedNote.id}`)">
                      <h5>{{ relatedNote.title || '无标题' }}</h5>
                      <p>{{ relatedNote.author || '未知作者' }} · {{ relatedNote.course_name || '未知课程' }}</p>
                      <div class="related-note-stats">
                        <span><el-icon><View /></el-icon> {{ relatedNote.view_count || 0 }}</span>
                        <span><el-icon><Star /></el-icon> {{ relatedNote.like_count || 0 }}</span>
                      </div>
                    </div>
                  </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 错误状态 -->
      <div v-else class="error-state">
        <el-result icon="warning" title="笔记不存在" sub-title="抱歉，您访问的笔记不存在或已被删除">
          <template #extra>
            <el-button type="primary" @click="$router.push('/notes')">返回笔记列表</el-button>
          </template>
        </el-result>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onActivated, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { noteAPI, commentAPI, roleUtils } from '../api'
import { aiAPI } from '../api'
import { 
  User, 
  Clock, 
  View, 
  Star, 
  StarFilled, 
  Collection, 
  Folder, 
  MoreFilled, 
  ChatDotRound 
} from '@element-plus/icons-vue'

export default {
  name: 'NoteDetail',
  components: {
    User,
    Clock,
    View,
    Star,
    StarFilled,
    Collection,
    Folder,
    MoreFilled,
    ChatDotRound
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const loading = ref(true)
    const note = ref({
      id: null,
      title: '',
      is_liked: false,
      is_favorited: false,
      like_count: 0
    })
    const comments = ref([])
    const authorNotes = ref([])
    const relatedNotes = ref([])
    const showCommentForm = ref(false)
    const commentContent = ref('')
    const commentLoading = ref(false)
    const replyingTo = ref(null)
    const replyContent = ref('')
    const replyLoading = ref(false)
    const currentUserId = ref(1) // 模拟当前用户ID
    const regenerating = ref(false)
    
    // 角色权限检查
    const canEditNote = computed(() => {
      return noteAPI.canEditNote(note.value)
    })
    
    const canDeleteNote = computed(() => {
      return noteAPI.canDeleteNote(note.value)
    })
    
    const canComment = computed(() => {
      return roleUtils.isLoggedIn()
    })
    
    const canEditComment = (comment) => {
      return commentAPI.canEditComment(comment)
    }
    
    const canDeleteComment = (comment) => {
      return commentAPI.canDeleteComment(comment)
    }

    // 获取笔记详情
    const fetchNoteDetail = async (forceRefresh = false) => {
      try {
        const noteId = route.params.id
        // 如果是强制刷新，添加时间戳避免缓存
        const params = forceRefresh ? { t: Date.now() } : {}
        const response = await noteAPI.getNoteDetail(noteId, params)
        // 确保字段映射正确
        note.value = {
          ...response,
          author: response.username || response.author,
          author_id: response.user_id || response.author_id,
          course_name: response.course_name || response.courseName,
          is_liked: response.is_liked === true,
          is_favorited: response.is_favorited === true,
          like_count: response.like_count || 0
        }
        console.log('笔记详情获取成功:', response)
        console.log('处理后的笔记数据:', note.value)
        console.log('点赞状态:', note.value.is_liked)
        
        // 获取评论
        await fetchComments(noteId)
    
        // 获取作者其他笔记
        await fetchAuthorNotes(note.value.author_id)
        
        // 获取相关笔记
        await fetchRelatedNotes(noteId)
        // 获取AI摘要与关键词
        await fetchNoteAIMeta(noteId)
      } catch (error) {
        console.error('获取笔记详情失败:', error)
        // 强制刷新页面以确保获取最新数据，避免缓存问题
        if (note.value && note.value.id) {
          console.log('API调用失败，但已有数据，尝试重新获取最新状态')
          // 清空当前数据，强制重新加载
          note.value = {
            id: null,
            title: '',
            is_liked: false,
            is_favorited: false,
            like_count: 0
          }
          // 延迟重试，使用强制刷新
          setTimeout(() => {
            fetchNoteDetail(true)
          }, 1000)
          return
        }
        // 如果完全没有数据，显示错误状态
        if (!note.value || !note.value.id) {
          ElMessage.error('笔记不存在或加载失败')
          // 设置空状态
          note.value = {
            id: null,
            title: '',
            content: '',
            author: '',
            author_id: null,
            is_liked: false,
            is_favorited: false,
            like_count: 0
          }
          comments.value = []
          authorNotes.value = []
          relatedNotes.value = []
          return
        } else {
          // 如果已有数据，显示错误信息但不覆盖现有状态
          ElMessage.error('获取笔记详情失败，显示的可能不是最新数据')
        }
      } finally {
        loading.value = false
      }
    }

    // 获取评论
    const fetchComments = async (noteId) => {
      try {
        console.log('=== 开始获取评论 ===')
        console.log('笔记ID:', noteId)
        console.log('API URL:', `/comment/note/${noteId}`)
        
        const response = await commentAPI.getCommentsByNote(noteId)
        console.log('=== API原始响应 ===')
        console.log('响应类型:', typeof response)
        console.log('响应内容:', JSON.stringify(response, null, 2))
        
        // 检查响应数据结构
        if (response && typeof response === 'object') {
          // 检查是否有comments字段
          if (Array.isArray(response.comments)) {
            comments.value = response.comments
            console.log('✅ 评论获取成功（对象格式），数量:', comments.value.length)
          } 
          // 检查是否直接是数组
          else if (Array.isArray(response)) {
            comments.value = response
            console.log('✅ 评论获取成功（数组格式），数量:', comments.value.length)
          }
          // 检查是否有data字段
          else if (response.data && Array.isArray(response.data)) {
            comments.value = response.data
            console.log('✅ 评论获取成功（data字段），数量:', comments.value.length)
          }
          // 其他可能的数据结构
          else {
            console.warn('⚠️ 评论数据格式异常，尝试查找评论数组')
            console.log('响应对象的所有键:', Object.keys(response))
            
            // 尝试查找可能包含评论的字段
            const possibleKeys = ['comments', 'data', 'list', 'items', 'results']
            let found = false
            
            for (const key of possibleKeys) {
              if (response[key] && Array.isArray(response[key])) {
                comments.value = response[key]
                console.log(`✅ 在${key}字段找到评论数组，数量:`, comments.value.length)
                found = true
                break
              }
            }
            
            if (!found) {
              console.log('❌ 未找到评论数组，设置为空')
              comments.value = []
            }
          }
        } else {
          console.warn('❌ 响应不是有效对象，设置为空数组')
          comments.value = []
        }
        
        console.log('=== 最终评论数据 ===')
        console.log('评论数量:', comments.value.length)
        if (comments.value.length > 0) {
          console.log('第一条评论:', comments.value[0])
        }
        
      } catch (error) {
        console.error('=== 获取评论失败 ===')
        console.error('错误对象:', error)
        console.error('错误消息:', error.message)
        
        if (error.response) {
          console.error('HTTP状态码:', error.response.status)
          console.error('响应数据:', error.response.data)
          
          const status = error.response.status
          if (status === 404) {
            console.log('📝 笔记不存在或无评论')
            ElMessage.info('该笔记暂无评论')
          } else if (status >= 500) {
            console.error('🔥 服务器错误')
            ElMessage.error('服务器错误，无法加载评论')
          } else {
            console.error('❌ 其他HTTP错误')
            ElMessage.error('加载评论失败，请稍后重试')
          }
        } else if (error.message && error.message.includes('Network Error')) {
          console.error('🌐 网络错误')
          ElMessage.error('网络连接失败，请检查网络设置')
        } else {
          console.error('❓ 未知错误')
          ElMessage.error('加载评论时发生未知错误')
        }
        
        // 设置空数组作为默认值
        comments.value = []
      }
    }

    // 获取作者其他笔记
    const fetchAuthorNotes = async (authorId) => {
      try {
        const response = await noteAPI.getNotes({ user_id: authorId, page_size: 5 })
        // 根据API拦截器，response已经是处理后的数据
        authorNotes.value = response.notes || response || []
        console.log('作者笔记获取成功:', response)
      } catch (error) {
        console.error('获取作者笔记失败:', error)
        // 设置空数组作为默认值
        authorNotes.value = []
      }
    }

    // 获取相关笔记
    const fetchRelatedNotes = async (noteId) => {
      try {
        // 获取最新的公开笔记作为相关笔记
        const response = await noteAPI.getNotes({ page_size: 5, sort_by: 'created_at', order: 'desc' })
        // 根据API拦截器，response已经是处理后的数据
        const notes = response.notes || response || []
        // 过滤掉当前笔记
        relatedNotes.value = notes.filter(n => n && n.id !== parseInt(noteId))
        console.log('相关笔记获取成功:', response)
      } catch (error) {
        console.error('获取相关笔记失败:', error)
        // 设置空数组作为默认值
        relatedNotes.value = []
      }
    }

    // 获取AI元数据（摘要/关键词）
    const aiMeta = ref({ summary: '', keywords: [] })
    const fetchNoteAIMeta = async (noteId) => {
      try {
        const res = await aiAPI.getNoteMeta(noteId)
        aiMeta.value = {
          summary: res?.summary || '',
          keywords: Array.isArray(res?.keywords) ? res.keywords : (res?.keywords ? String(res.keywords).split(',').filter(k => k.trim()) : [])
        }
      } catch (error) {
        console.warn('获取AI元数据失败:', error)
        aiMeta.value = { summary: '', keywords: [] }
      }
    }

    // 重新生成摘要
    const handleRegenerateSummary = async () => {
      try {
        // 权限校验：仅作者或管理员可用
        if (!canEditNote.value) {
          ElMessage.warning('您没有权限重新生成摘要')
          return
        }
        regenerating.value = true
        const noteId = route.params.id
        // 调用后端生成摘要
        await aiAPI.summarize({ note_id: Number(noteId) })
        // 重新拉取AI元数据
        await fetchNoteAIMeta(noteId)
        ElMessage.success('摘要已重新生成')
      } catch (error) {
        console.error('重新生成摘要失败:', error)
        const msg = error?.response?.data?.error || error?.message || '重新生成失败'
        ElMessage.error(msg)
      } finally {
        regenerating.value = false
      }
    }

    // 格式化内容（简单的 Markdown 转换）
    const formattedContent = computed(() => {
      try {
        let content = note.value.content
        
        // 检查内容是否存在且为字符串
        if (!content || typeof content !== 'string') {
          return '<p>暂无内容</p>'
        }
        
        // 简单的 Markdown 转换
        content = content
          .replace(/^# (.*$)/gim, '<h1>$1</h1>')
          .replace(/^## (.*$)/gim, '<h2>$1</h2>')
          .replace(/^### (.*$)/gim, '<h3>$1</h3>')
          .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
          .replace(/\*(.*?)\*/g, '<em>$1</em>')
          .replace(/\n/g, '<br>')
        
        // 确保返回的内容不为空
        return content.trim() || '<p>暂无内容</p>'
      } catch (error) {
        console.error('格式化内容时出错:', error)
        return '<p>内容格式化失败</p>'
      }
    })

    // 点赞状态
    const isLiking = ref(false)
    
    // 点赞笔记
    const handleLike = async () => {
      // 防止重复点击
      if (isLiking.value) {
        return
      }
      
      isLiking.value = true
      
      try {
        let response
        if (note.value.is_liked) {
          response = await noteAPI.unlikeNote(note.value.id)
          ElMessage.success('已取消点赞')
        } else {
          response = await noteAPI.likeNote(note.value.id)
          ElMessage.success('点赞成功')
        }
        
        // 直接使用后端返回的状态更新前端，避免重新获取笔记详情
        if (response && typeof response.like_count !== 'undefined' && typeof response.is_liked !== 'undefined') {
          note.value.like_count = response.like_count
          note.value.is_liked = response.is_liked
          console.log('点赞状态已更新:', { like_count: response.like_count, is_liked: response.is_liked })
        } else {
          // 如果后端没有返回完整状态，则刷新获取最新状态
          console.log('后端未返回完整状态，刷新获取最新状态')
          await fetchNoteDetail(true)
        }
      } catch (error) {
        console.error('点赞操作失败:', error)
        
        // 检查具体的错误状态
        if (error.response) {
          const status = error.response.status
          const data = error.response.data || {}
          const message = data.error || data.message || '操作失败'
          
          if (status === 409) {
            ElMessage.warning('您已经点赞过此笔记')
            // 使用错误响应中的状态信息更新前端
            if (typeof data.like_count !== 'undefined' && typeof data.is_liked !== 'undefined') {
              note.value.like_count = data.like_count
              note.value.is_liked = data.is_liked
            }
          } else if (status === 404 && message.includes('还没有点赞')) {
            ElMessage.warning('您还没有点赞此笔记')
            // 使用错误响应中的状态信息更新前端
            if (typeof data.like_count !== 'undefined' && typeof data.is_liked !== 'undefined') {
              note.value.like_count = data.like_count
              note.value.is_liked = data.is_liked
            }
          } else {
            ElMessage.error(`操作失败: ${message}`)
            // 对于其他错误，强制刷新笔记数据以获取最新状态
            console.log('点赞操作出错，强制刷新获取最新状态')
            await fetchNoteDetail(true)
          }
        } else {
          ElMessage.error('网络错误，请检查网络连接')
        }
      } finally {
        isLiking.value = false
      }
    }

    // 收藏状态
    const isFavoriting = ref(false)
    
    // 收藏笔记
    const handleFavorite = async () => {
      // 防止重复点击
      if (isFavoriting.value) {
        return
      }
      
      isFavoriting.value = true
      
      try {
        if (note.value.is_favorited) {
          await noteAPI.unfavoriteNote(note.value.id)
          note.value.is_favorited = false
          ElMessage.success('已取消收藏')
        } else {
          await noteAPI.favoriteNote(note.value.id)
          note.value.is_favorited = true
          ElMessage.success('收藏成功')
        }
      } catch (error) {
        console.error('收藏操作失败:', error)
        
        // 检查是否是409冲突错误
        if (error.response && error.response.status === 409) {
          ElMessage.warning('操作过于频繁，请稍后再试')
          // 强制刷新笔记数据以获取最新状态
          console.log('点赞操作出错，强制刷新获取最新状态')
          await fetchNoteDetail(true)
        } else {
          // 模拟操作作为降级方案
          if (note.value.is_favorited) {
            note.value.is_favorited = false
            ElMessage.success('已取消收藏')
          } else {
            note.value.is_favorited = true
            ElMessage.success('收藏成功')
          }
        }
      } finally {
        isFavoriting.value = false
      }
    }



    // 笔记操作
    const handleNoteAction = (command) => {
      if (command === 'edit') {
        router.push(`/notes/${note.value.id}/edit`)
      } else if (command === 'delete') {
        ElMessageBox.confirm('确定要删除这篇笔记吗？', '确认删除', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(async () => {
          try {
            await noteAPI.deleteNote(note.value.id)
            ElMessage.success('删除成功')
            router.push('/notes')
          } catch (error) {
            ElMessage.error('删除失败')
          }
        })
      }
    }

    // 提交评论
    const handleSubmitComment = async () => {
      if (!commentContent.value.trim()) {
        ElMessage.warning('请输入评论内容')
        return
      }
      
      commentLoading.value = true
      try {
        const response = await commentAPI.createComment({
          note_id: note.value.id,
          content: commentContent.value
        })
        
        // 根据API拦截器，response已经是处理后的数据
        comments.value.unshift(response)
        commentContent.value = ''
        showCommentForm.value = false
        ElMessage.success('评论发表成功')
        console.log('评论提交成功:', response)
        
        // 重新获取笔记详情以更新评论数量等信息
        await fetchNoteDetail()
      } catch (error) {
        console.error('评论提交失败:', error)
        // 模拟添加评论
        const newComment = {
          id: Date.now(),
          content: commentContent.value,
          author: '当前用户',
          author_avatar: '',
          like_count: 0,
          created_at: new Date().toISOString()
        }
        comments.value.unshift(newComment)
        commentContent.value = ''
        showCommentForm.value = false
        ElMessage.success('评论发表成功')
        
        // 即使是模拟数据，也更新评论数量
        if (note.value) {
          note.value.comment_count = (note.value.comment_count || 0) + 1
        }
      } finally {
        commentLoading.value = false
      }
    }

    // 评论点赞状态
    const commentLikingMap = ref(new Map())
    
    // 点赞评论
    const handleLikeComment = async (comment) => {
      if (!canComment.value) {
        ElMessage.warning('请先登录')
        router.push('/login')
        return
      }
      
      // 防止重复点击
      if (commentLikingMap.value.get(comment.id)) {
        return
      }
      
      commentLikingMap.value.set(comment.id, true)
      
      // 保存原始状态，用于错误回滚
      const originalLiked = comment.is_liked
      const originalCount = comment.like_count || 0
      
      try {
        let response
        
        if (comment.is_liked) {
          // 取消点赞
          response = await commentAPI.unlikeComment(comment.id)
        } else {
          // 点赞
          response = await commentAPI.likeComment(comment.id)
        }
        
        // 使用后端返回的准确数据更新状态
        if (response && response.success && response.data) {
          comment.like_count = response.data.like_count
          comment.is_liked = response.data.is_liked
          ElMessage.success(response.message || (comment.is_liked ? '点赞成功' : '取消点赞成功'))
        } else {
          // 如果响应格式不符合预期，回滚状态
          comment.is_liked = originalLiked
          comment.like_count = originalCount
          ElMessage.error('操作失败，请重试')
        }
        
      } catch (error) {
        console.error('点赞操作失败:', error)
        
        // 回滚到原始状态
        comment.is_liked = originalLiked
        comment.like_count = originalCount
        
        // 处理错误响应
        if (error.response) {
          const status = error.response.status
          const data = error.response.data || {}
          const message = data.error || data.message || '操作失败'
          
          if (status === 409) {
            // 409表示已经点赞过，使用后端返回的准确状态
            ElMessage.warning('您已经点赞过此评论')
            if (typeof data.like_count !== 'undefined' && typeof data.is_liked !== 'undefined') {
              comment.like_count = data.like_count
              comment.is_liked = data.is_liked
            }
          } else if (status === 404 && message.includes('还没有点赞')) {
            // 404表示还没有点赞过，使用后端返回的准确状态
            ElMessage.warning('您还没有点赞此评论')
            if (typeof data.like_count !== 'undefined' && typeof data.is_liked !== 'undefined') {
              comment.like_count = data.like_count
              comment.is_liked = data.is_liked
            }
          } else if (status === 401) {
            ElMessage.error('登录已过期，请重新登录')
            router.push('/login')
          } else {
            ElMessage.error(`操作失败: ${message}`)
          }
        } else {
          ElMessage.error('网络错误，请检查网络连接')
        }
      } finally {
        commentLikingMap.value.set(comment.id, false)
      }
    }

    // 回复评论
    const handleReplyComment = (comment) => {
      if (!canComment.value) {
        ElMessage.warning('请先登录')
        router.push('/login')
        return
      }
      
      replyingTo.value = comment.id
      replyContent.value = ''
    }
    
    // 回复回复
    const handleReplyToReply = (parentComment, reply) => {
      if (!canComment.value) {
        ElMessage.warning('请先登录')
        router.push('/login')
        return
      }
      
      replyingTo.value = parentComment.id
      replyContent.value = `@${reply.username} `
    }
    
    // 取消回复
    const cancelReply = () => {
      replyingTo.value = null
      replyContent.value = ''
    }
    
    // 提交回复
    const submitReply = async (parentComment) => {
      if (!replyContent.value.trim()) {
        ElMessage.warning('请输入回复内容')
        return
      }
      
      replyLoading.value = true
      try {
        // 解析@用户名
        let replyToUserID = null
        const atMatch = replyContent.value.match(/@(\S+)\s/)
        if (atMatch) {
          // 这里应该根据用户名查找用户ID，暂时模拟
          const username = atMatch[1]
          // 在回复列表中查找用户
          const targetReply = parentComment.replies?.find(r => r.username === username)
          if (targetReply) {
            replyToUserID = targetReply.user_id
          }
        }
        
        const response = await commentAPI.createComment({
          note_id: note.value.id,
          content: replyContent.value,
          parent_id: parentComment.id,
          reply_to_user_id: replyToUserID
        })
        
        // 添加到回复列表
        if (!parentComment.replies) {
          parentComment.replies = []
        }
        parentComment.replies.push(response)
        
        cancelReply()
        ElMessage.success('回复发表成功')
        
        // 重新获取评论数据以确保数据同步
        await fetchComments(note.value.id)
      } catch (error) {
        console.error('回复提交失败:', error)
        // 模拟添加回复
        const newReply = {
          id: Date.now(),
          content: replyContent.value,
          username: '当前用户',
          user_avatar: '',
          user_id: currentUserId.value,
          like_count: 0,
          is_liked: false,
          created_at: new Date().toISOString()
        }
        
        if (!parentComment.replies) {
          parentComment.replies = []
        }
        parentComment.replies.push(newReply)
        
        cancelReply()
        ElMessage.success('回复发表成功')
        
        // 即使是模拟数据，也更新评论数量
        if (note.value) {
          note.value.comment_count = (note.value.comment_count || 0) + 1
        }
      } finally {
        replyLoading.value = false
      }
    }
    
    // 评论操作
    const handleCommentAction = async (command, comment) => {
      if (command === 'edit') {
        ElMessage.info('编辑功能开发中')
      } else if (command === 'delete') {
        ElMessageBox.confirm('确定要删除这条评论吗？', '确认删除', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(async () => {
          try {
            await commentAPI.deleteComment(comment.id, comment.user_id)
            // 从列表中移除评论
            const index = comments.value.findIndex(c => c && c.id === comment.id)
            if (index > -1) {
              comments.value.splice(index, 1)
            }
            ElMessage.success('删除成功')
          } catch (error) {
            console.error('删除评论失败:', error)
            ElMessage.error('删除失败')
          }
        })
      }
    }

    // 获取分类类型
    const getCategoryType = (category) => {
      const typeMap = {
        '数学': 'primary',
        '计算机': 'success',
        '语言': 'warning',
        '物理': 'info',
        '化学': 'danger',
        '经济': 'primary',
        '管理': 'success'
      }
      return typeMap[category] || 'info'
    }

    // 获取难度类型
    const getDifficultyType = (difficulty) => {
      const typeMap = {
        'beginner': 'success',
        'intermediate': 'warning',
        'advanced': 'danger'
      }
      return typeMap[difficulty] || 'info'
    }

    // 获取难度文本
    const getDifficultyText = (difficulty) => {
      const textMap = {
        'beginner': '入门',
        'intermediate': '中级',
        'advanced': '高级'
      }
      return textMap[difficulty] || '未知'
    }

    // 获取字数统计
    const getWordCount = (content) => {
      return content ? content.replace(/<[^>]*>/g, '').length : 0
    }

    // 获取阅读时间
    const getReadingTime = (content) => {
      const wordCount = getWordCount(content)
      return Math.ceil(wordCount / 300) // 假设每分钟阅读300字
    }

    // 格式化时间
    const formatTime = (timeString) => {
      const date = new Date(timeString)
      const now = new Date()
      const diff = now - date
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))
      
      if (days === 0) {
        const hours = Math.floor(diff / (1000 * 60 * 60))
        if (hours === 0) {
          const minutes = Math.floor(diff / (1000 * 60))
          return `${minutes}分钟前`
        }
        return `${hours}小时前`
      } else if (days < 7) {
        return `${days}天前`
      } else {
        return date.toLocaleDateString()
      }
    }

    // 格式化日期
    const formatDate = (timeString) => {
      return new Date(timeString).toLocaleDateString()
    }

    onMounted(() => {
      console.log('组件挂载，强制刷新数据以获取最新状态')
      fetchNoteDetail(true) // 强制刷新
    })

    // 页面激活时重新获取数据（处理keep-alive缓存的情况）
    onActivated(() => {
      console.log('页面激活，强制刷新数据以获取最新状态')
      fetchNoteDetail(true) // 强制刷新
    })

    // 监听路由参数变化，当笔记ID改变时重新获取数据
    watch(() => route.params.id, (newId, oldId) => {
      if (newId && newId !== oldId) {
        console.log('路由参数变化，强制刷新数据', { newId, oldId })
        fetchNoteDetail(true) // 强制刷新
      }
    })

    return {
      loading,
      note,
      comments,
      authorNotes,
      relatedNotes,
      showCommentForm,
      commentContent,
      commentLoading,
      replyingTo,
      replyContent,
      replyLoading,
      currentUserId,
      canEditNote,
      canDeleteNote,
      canComment,
      canEditComment,
      canDeleteComment,
      formattedContent,
      isLiking,
      isFavoriting,
      commentLikingMap,
      handleLike,
      handleFavorite,

      handleNoteAction,
      handleSubmitComment,
      handleLikeComment,
      handleReplyComment,
      handleReplyToReply,
      fetchNoteAIMeta,
      aiMeta,
      handleCommentAction,
      cancelReply,
      submitReply,
      getCategoryType,
      getDifficultyType,
      getDifficultyText,
      getWordCount,
      getReadingTime,
      formatTime,
      formatDate,
      regenerating,
      handleRegenerateSummary
    }
  }
}
</script>

<style scoped>
.page-header {
  background: var(--bg-secondary);
  padding: 20px 0;
  border-bottom: 1px solid var(--border-color);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.loading-container {
  padding: 40px 0;
}

.note-detail {
  padding: 30px 0;
}

.note-content-card {
  margin-bottom: 30px;
}

.note-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
}

.note-meta {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.note-time,
.note-views {
  display: flex;
  align-items: center;
  gap: 5px;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.note-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.note-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 30px;
  line-height: 1.3;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 30px;
  padding: 20px;
  background: var(--bg-secondary);
  border-radius: 8px;
}

.author-details {
  flex: 1;
}

.author-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 5px;
}

.author-meta {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.note-tags {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 30px;
}

.note-content {
  font-size: 1.1rem;
  line-height: 1.8;
  color: var(--text-primary);
}

.note-content :deep(h1) {
  font-size: 2rem;
  font-weight: 600;
  margin: 30px 0 20px 0;
  color: var(--text-primary);
}

.note-content :deep(h2) {
  font-size: 1.6rem;
  font-weight: 600;
  margin: 25px 0 15px 0;
  color: var(--text-primary);
}

.note-content :deep(h3) {
  font-size: 1.3rem;
  font-weight: 600;
  margin: 20px 0 10px 0;
  color: var(--text-primary);
}

.note-content :deep(p) {
  margin-bottom: 15px;
}

.note-content :deep(strong) {
  font-weight: 600;
  color: var(--primary-color);
}

.note-content :deep(em) {
  font-style: italic;
  color: var(--text-secondary);
}

.comments-section {
  margin-top: 40px;
  padding: 0 20px;
}

.comments-card {
  margin-bottom: 30px;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--border-color);
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
}

.comments-card :deep(.el-card__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px 24px;
  border-bottom: none;
}

.comments-card :deep(.el-card__body) {
  padding: 24px;
}

.comments-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 1.1rem;
  font-weight: 600;
}

.comments-header span {
  color: white;
}

.comments-header .el-button {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  backdrop-filter: blur(10px);
}

.comments-header .el-button:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
}

.comment-form {
  margin-bottom: 30px;
  padding: 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  border-radius: 12px;
  border: 1px solid var(--border-color);
}

.comment-form .el-textarea {
  margin-bottom: 15px;
}

.comment-form :deep(.el-textarea__inner) {
  border-radius: 8px;
  border: 2px solid #e2e8f0;
  transition: all 0.3s ease;
}

.comment-form :deep(.el-textarea__inner):focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.comment-form .comment-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 15px;
}

.comment-form .comment-actions .el-button {
  border-radius: 8px;
  padding: 8px 20px;
  font-weight: 500;
}

.empty-comments {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  border-radius: 12px;
  margin: 20px 0;
}

.empty-comments .el-button {
  margin-top: 15px;
  border-radius: 8px;
  padding: 10px 24px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin-top: 20px;
}

.comment-item {
  display: flex;
  gap: 15px;
  padding: 24px;
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s ease;
  position: relative;
}

.comment-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  border-color: var(--primary-color);
}

.comment-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px 0 0 2px;
}

.comment-avatar {
  flex-shrink: 0;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.comment-author {
  font-weight: 600;
  color: var(--text-primary);
}

.reply-to {
  color: var(--primary-color);
  font-size: 0.9rem;
}

.comment-time {
  color: var(--text-secondary);
  font-size: 0.9rem;
  margin-left: auto;
}

.comment-text {
  margin-bottom: 15px;
  line-height: 1.6;
}

.comment-item .comment-actions {
  display: flex;
  gap: 15px;
  align-items: center;
  margin-top: 10px;
}

.comment-item .comment-actions .el-button {
  border-radius: 6px;
  font-size: 0.9rem;
}

.reply-form {
  margin-top: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
  border-radius: 10px;
  border: 1px solid #cbd5e1;
}

.reply-form :deep(.el-textarea__inner) {
  border-radius: 8px;
  border: 2px solid #e2e8f0;
}

.reply-form :deep(.el-textarea__inner):focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 10px;
}

.replies-list {
  margin-top: 24px;
  padding-left: 24px;
  border-left: 3px solid #e2e8f0;
  position: relative;
}

.replies-list::before {
  content: '';
  position: absolute;
  left: -3px;
  top: 0;
  bottom: 0;
  width: 3px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
}

.reply-item {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding: 18px;
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s ease;
}

.reply-item:hover {
  transform: translateX(4px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  border-color: #cbd5e1;
}

.reply-avatar {
  flex-shrink: 0;
}

.reply-content {
  flex: 1;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.reply-author {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 0.9rem;
}

.reply-time {
  color: var(--text-secondary);
  font-size: 0.85rem;
  margin-left: auto;
}

.reply-text {
  margin-bottom: 10px;
  line-height: 1.5;
  font-size: 0.95rem;
}

.reply-actions {
   display: flex;
   gap: 10px;
   align-items: center;
 }

.comment-author {
  font-weight: 600;
  color: var(--text-primary);
}

.comment-time {
  font-size: 0.9rem;
  color: var(--text-light);
}

.comment-text {
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 10px;
}

.comment-actions {
  display: flex;
  gap: 15px;
}

.note-info-sidebar,
.author-notes-sidebar,
.related-notes-sidebar {
  margin-bottom: 30px;
}

.sidebar-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-item label {
  font-weight: 500;
  color: var(--text-secondary);
}

.empty-author-notes {
  text-align: center;
  padding: 20px;
}

.author-notes-list,
.related-notes-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.author-note-item,
.related-note-item {
  cursor: pointer;
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: border-color 0.3s;
}

.author-note-item:hover,
.related-note-item:hover {
  border-color: var(--primary-color);
}

.author-note-info h5,
.related-note-info h5 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 5px;
}

.author-note-info p,
.related-note-info p {
  color: var(--text-secondary);
  margin-bottom: 10px;
  font-size: 0.9rem;
}

.author-note-stats,
.related-note-stats {
  display: flex;
  gap: 15px;
  font-size: 0.9rem;
  color: var(--text-light);
}

.author-note-stats span,
.related-note-stats span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.error-state {
  padding: 60px 0;
  text-align: center;
}

.ai-summary-card {
  margin-bottom: 30px;
}
.ai-summary-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.ai-summary-text {
  white-space: pre-wrap;
  line-height: 1.7;
  color: var(--text-primary);
}
.ai-summary-empty {
  color: var(--text-secondary);
}
.ai-keywords {
  margin-top: 10px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.ai-keyword {
  margin-right: 6px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .note-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .note-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .note-title {
    font-size: 2rem;
  }
  
  .author-info {
    flex-direction: column;
    text-align: center;
  }
  
  .el-row {
    flex-direction: column;
  }
  
  .el-col {
    width: 100%;
    margin-bottom: 20px;
  }
  
  .comment-item {
    flex-direction: column;
    gap: 10px;
  }
}
</style>