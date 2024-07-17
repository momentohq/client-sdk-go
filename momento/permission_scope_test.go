package momento_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/momentohq/client-sdk-go/internal"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
)

var _ = Describe("auth auth-client", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCaches()

		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	Describe("PermissionScope", func() {
		It("should support assignment from PredefinedScope and AllDataReadWrite", func() {
			var scope PermissionScope

			scope = internal.InternalSuperUserPermissions{}
			Expect(scope).To(Equal(internal.InternalSuperUserPermissions{}))

			scope = AllDataReadWrite
			Expect(scope).To(Equal(&Permissions{
				Permissions: []Permission{
					TopicPermission{Topic: AllTopics{}, Cache: AllCaches{}, Role: PublishSubscribe},
					CachePermission{Cache: AllCaches{}, Role: ReadWrite},
				},
			}))
		})

		Describe("should support assignment from Permissions literal", func() {
			var scope PermissionScope

			It("using cache name in a CachePermission", func() {
				scope = Permissions{
					Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: "my-cache"}, Role: ReadWrite},
					},
				}
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: "my-cache"}, Role: ReadWrite},
					},
				}))
			})

			It("using AllCaches in a CachePermission", func() {
				scope = Permissions{
					Permissions: []Permission{
						CachePermission{Cache: AllCaches{}, Role: ReadWrite},
					},
				}
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: AllCaches{}, Role: ReadWrite},
					},
				}))
			})

			It("using cache name and topic name in a TopicPermission", func() {
				scope = Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: TopicName{Name: "my-topic"}, Role: PublishSubscribe},
					},
				}
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: TopicName{Name: "my-topic"}, Role: PublishSubscribe},
					},
				}))
			})

			It("using cache name and AllTopics in a TopicPermission", func() {
				scope = Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: AllTopics{}, Role: PublishSubscribe},
					},
				}
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: AllTopics{}, Role: PublishSubscribe},
					},
				}))
			})

			It("using AllCaches and topic name in a TopicPermission", func() {
				scope = Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: AllCaches{}, Topic: TopicName{Name: "my-topic"}, Role: PublishSubscribe},
					},
				}
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: AllCaches{}, Topic: TopicName{Name: "my-topic"}, Role: PublishSubscribe},
					},
				}))
			})

			It("using AllCaches and AllTopics in a TopicPermission", func() {
				scope = Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: AllCaches{}, Topic: AllTopics{}, Role: PublishSubscribe},
					},
				}
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: AllCaches{}, Topic: AllTopics{}, Role: PublishSubscribe},
					},
				}))
			})

			It("mixing cache and topic permissions", func() {
				scope = Permissions{
					Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: "my-cache"}, Role: ReadWrite},
						TopicPermission{Cache: CacheName{Name: "another-cache"}, Topic: TopicName{Name: "topic1"}, Role: PublishSubscribe},
						TopicPermission{Cache: CacheName{Name: "another-cache"}, Topic: TopicName{Name: "topic2"}, Role: PublishSubscribe},
					},
				}
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: "my-cache"}, Role: ReadWrite},
						TopicPermission{Cache: CacheName{Name: "another-cache"}, Topic: TopicName{Name: "topic1"}, Role: PublishSubscribe},
						TopicPermission{Cache: CacheName{Name: "another-cache"}, Topic: TopicName{Name: "topic2"}, Role: PublishSubscribe},
					},
				}))
			})
		})

		Describe("should support assignment from PermissionScope factory functions", func() {
			It("CacheReadWrite", func() {
				// Specific cache name
				scope := CacheReadWrite(CacheName{Name: "my-cache"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: "my-cache"}, Role: ReadWrite},
					},
				}))

				// AllCaches
				scope = CacheReadWrite(AllCaches{})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: AllCaches{}, Role: ReadWrite},
					},
				}))
			})

			It("CacheReadOnly", func() {
				// Specific cache name
				scope := CacheReadOnly(CacheName{Name: "my-cache"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: "my-cache"}, Role: ReadOnly},
					},
				}))

				// AllCaches
				scope = CacheReadOnly(AllCaches{})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: AllCaches{}, Role: ReadOnly},
					},
				}))
			})

			It("CacheWriteOnly", func() {
				// Specific cache name
				scope := CacheWriteOnly(CacheName{Name: "my-cache"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: "my-cache"}, Role: WriteOnly},
					},
				}))

				// AllCaches
				scope = CacheWriteOnly(AllCaches{})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						CachePermission{Cache: AllCaches{}, Role: WriteOnly},
					},
				}))
			})

			It("TopicSubscribeOnly", func() {
				// Specific cache and topic
				scope := TopicSubscribeOnly(CacheName{Name: "my-cache"}, TopicName{Name: "my-topic"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: TopicName{Name: "my-topic"}, Role: SubscribeOnly},
					},
				}))

				// Specific cache and AllTopics
				scope = TopicSubscribeOnly(CacheName{Name: "my-cache"}, AllTopics{})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: AllTopics{}, Role: SubscribeOnly},
					},
				}))

				// AllCaches and specific topic
				scope = TopicSubscribeOnly(AllCaches{}, TopicName{Name: "my-topic"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: AllCaches{}, Topic: TopicName{Name: "my-topic"}, Role: SubscribeOnly},
					},
				}))
			})

			It("TopicPublishSubscribe", func() {
				// Specific cache and topic
				scope := TopicPublishSubscribe(CacheName{Name: "my-cache"}, TopicName{Name: "my-topic"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: TopicName{Name: "my-topic"}, Role: PublishSubscribe},
					},
				}))

				// Specific cache and AllTopics
				scope = TopicPublishSubscribe(CacheName{Name: "my-cache"}, AllTopics{})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: AllTopics{}, Role: PublishSubscribe},
					},
				}))

				// AllCaches and specific topic
				scope = TopicPublishSubscribe(AllCaches{}, TopicName{Name: "my-topic"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: AllCaches{}, Topic: TopicName{Name: "my-topic"}, Role: PublishSubscribe},
					},
				}))
			})

			It("TopicPublishOnly", func() {
				// Specific cache and topic
				scope := TopicPublishOnly(CacheName{Name: "my-cache"}, TopicName{Name: "my-topic"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: TopicName{Name: "my-topic"}, Role: PublishOnly},
					},
				}))

				// Specific cache and AllTopics
				scope = TopicPublishOnly(CacheName{Name: "my-cache"}, AllTopics{})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "my-cache"}, Topic: AllTopics{}, Role: PublishOnly},
					},
				}))

				// AllCaches and specific topic
				scope = TopicPublishOnly(AllCaches{}, TopicName{Name: "my-topic"})
				Expect(scope).To(Equal(Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: AllCaches{}, Topic: TopicName{Name: "my-topic"}, Role: PublishOnly},
					},
				}))
			})
		})
	})

})
